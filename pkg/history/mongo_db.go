package history

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/ip2country"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
)

type mongoDBClient struct {
	coll       *mongo.Collection
	client     *mongo.Client
	validator  *validator.Validator
	ip2country *ip2country.ClientWithCache
}

type MongoDBConfig struct {
	URL  string
	Name string
}

func NewMongoDBClient(cfg *MongoDBConfig, validatorEngine *validator.Validator, cache *ip2country.ClientWithCache) (Client, error) {
	mClient := &mongoDBClient{
		validator:  validatorEngine,
		ip2country: cache,
	}
	var (
		err error
	)

	mClient.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URL))
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	err = mClient.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	mClient.coll = mClient.client.Database(cfg.Name).Collection(SpinsCollectionName)

	// Get existing indexes
	ctx := context.Background()
	cursor, err := mClient.coll.Indexes().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var indexes []mongo.IndexModel
	if err = cursor.All(ctx, &indexes); err != nil {
		log.Fatal(err)
	}

	mClient.createIndexIfNotExists(ctx, indexes, "created_at", 1) // For B-Tree index

	mClient.createIndexIfNotExists(ctx, indexes, "currency", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "external_user_id", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "game", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "game_id", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "host", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "integrator", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "internal_user_id", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "provider", "hashed")
	mClient.createIndexIfNotExists(ctx, indexes, "session_token", "hashed")

	return mClient, nil
}

func existsIndex(indexes []mongo.IndexModel, key string, value interface{}) bool {
	for _, index := range indexes {
		if index.Keys != nil {
			for _, elem := range index.Keys.(bson.D) {
				if elem.Key == key && elem.Value == value {
					return true
				}
			}
		}
	}
	return false
}

func (m *mongoDBClient) createIndexIfNotExists(ctx context.Context, indexes []mongo.IndexModel, key string, value interface{}) {
	if !existsIndex(indexes, key, value) {
		indexModel := mongo.IndexModel{
			Keys: bson.D{{Key: key, Value: value}},
		}
		_, err := m.coll.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *mongoDBClient) Create(ctx context.Context, record *SpinIn) error {
	spin, err := spinIn2Spin(record)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	if err := spin.validateSpin(m.validator); err != nil {
		return err
	}

	country, err := m.ip2country.Get(spin.ClientIP)

	if err != nil {
		zap.S().Error(err)
	} else {
		spin.Country = &country
	}

	spin.Day = spin.CreatedAt

	_, err = m.coll.InsertOne(ctx, spin)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDBClient) Update(ctx context.Context, record *SpinIn) error {
	spin, err := spinIn2Spin(record)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	if err := spin.validateSpin(m.validator); err != nil {
		return err
	}

	country, err := m.ip2country.Get(spin.ClientIP)

	if err != nil {
		zap.S().Error(err)
	} else {
		spin.Country = &country
	}

	spin.Day = spin.CreatedAt

	spinByte, err := bson.Marshal(spin)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	var update bson.M
	err = bson.Unmarshal(spinByte, &update)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	err = m.coll.FindOneAndUpdate(ctx,
		bson.D{{"id", spin.ID}},
		bson.D{{Key: "$set", Value: update}}, options.FindOneAndUpdate().SetUpsert(true)).Decode(&spin)
	if err != nil {
		return err
	}
	return nil

}

func (m *mongoDBClient) Pagination(ctx context.Context, internalUserID uuid.UUID, game string, count int, page int) (p *GetSpinPaginationOut, err error) {
	var (
		records  []*Spin
		countI64 = int64(count)
		skip     = int64(page*count - count)
	)

	filter := bson.D{{"internal_user_id", internalUserID.String()},
		{"game", game},
		{"is_shown", true},
	}

	total, err := m.coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	cursor, err := m.coll.Find(ctx, filter,
		options.Find().SetSort(bson.D{{"created_at", -1}}), &options.FindOptions{Limit: &countI64, Skip: &skip})

	if err = cursor.All(ctx, &records); err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return &GetSpinPaginationOut{
		Items: lo.Map(records, func(item *Spin, index int) *SpinOut {
			return item.ToAPIResponse()
		}),
		Page:  uint64(page),
		Limit: uint64(count),
		Total: uint64(total),
	}, nil
}

func (m *mongoDBClient) LastRecords(ctx context.Context, internalUserID uuid.UUID, game string) ([]*SpinOut, error) {
	var records []*Spin

	cursor, err := m.coll.Find(ctx,
		bson.D{{"internal_user_id", internalUserID.String()},
			{"game", game},
			{"is_shown", false},
		},
		options.Find().SetSort(bson.D{{"created_at", -1}}))

	if err = cursor.All(ctx, &records); err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return lo.Map(records, func(item *Spin, index int) *SpinOut {
		return item.ToAPIResponse()
	}), nil
}

func (m *mongoDBClient) LastRecord(ctx context.Context, internalUserID uuid.UUID, game string) (*SpinOut, error) {
	return m.getBy(ctx,
		bson.D{{"internal_user_id", internalUserID.String()}, {"game", game}},
		options.FindOne().SetSort(bson.D{{"created_at", -1}}))
}

func (m *mongoDBClient) LastRecordByWager(ctx context.Context, internalUserID uuid.UUID, game string, wager uint64) (*SpinOut, error) {
	return m.getBy(ctx,
		bson.D{{"internal_user_id", internalUserID.String()},
			{"game", game},
			{"wager", wager},
		},
		options.FindOne().SetSort(bson.D{{"created_at", -1}}))
}

func (m *mongoDBClient) GetByID(ctx context.Context, id uuid.UUID) (*SpinOut, error) {
	return m.getBy(ctx, bson.D{{"id", id.String()}})
}

func (m *mongoDBClient) getBy(ctx context.Context, filter bson.D, opts ...*options.FindOneOptions) (*SpinOut, error) {
	var (
		spin *Spin
		err  error
	)

	err = m.coll.FindOne(ctx, filter, opts...).Decode(&spin)
	if err != nil {
		zap.S().Error(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSpinNotFound
		}
		return nil, err
	}

	return spin.ToAPIResponse(), nil
}
