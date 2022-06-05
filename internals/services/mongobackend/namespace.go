package mongobackend

import (
	"context"
	"log"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NamespaceService struct {
	BaseService
	varCollection *mongo.Collection
}

func NewNamespaceService(config *config.Config) *NamespaceService {
	nssvc := &NamespaceService{}
	nssvc.Connect(config)
	nssvc.collection = nssvc.client.Collection("namespaces")
	nssvc.varCollection = nssvc.client.Collection("namespacevars")
	return nssvc
}

func (nsSvc *NamespaceService) CreateNamespace(orgId string, name string, schemaId string, schemaVersion string, vars map[string]string) (*services.Namespace, error) {
	ns := &services.Namespace{
		Id:             uuid.NewString(),
		Name:           name,
		OrganizationId: orgId,
		SchemaId:       schemaId,
		SchemaVersion:  schemaVersion,
	}

	exists := nsSvc.ExistsByName(orgId, name)
	if exists {
		return nil, services.ErrNamespaceAlreadyExists
	}
	ctx := context.Background()
	_, err := nsSvc.collection.InsertOne(ctx, ns)
	if err != nil {
		log.Printf("Failed to create namespace %v", err)
		return nil, err
	}

	return ns, nil
}

func (nsSvc *NamespaceService) ExistsByName(orgId string, nsName string) bool {

	filter := bson.M{
		"name":           nsName,
		"organizationid": orgId,
	}

	ctx := context.Background()
	result := nsSvc.collection.FindOne(ctx, filter)

	return result.Err() != mongo.ErrNoDocuments
}

func (nsSvc *NamespaceService) GetNamespaceById(orgId string, nsId string) (*services.Namespace, error) {
	filter := bson.M{
		"id":             nsId,
		"organizationid": orgId,
	}

	ctx := context.Background()
	result := nsSvc.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		log.Printf("failed to get namespace by ID: %v", result.Err())
		return nil, result.Err()
	}

	var ns services.Namespace
	err := result.Decode(&ns)
	if err != nil {
		log.Printf("failed to decode while get namespace by ID: %v", err)
		return nil, err
	}
	return &ns, nil
}

func (nsSvc *NamespaceService) ListNamespaces(orgId string) ([]*services.Namespace, error) {
	filter := bson.M{
		"organizationid": orgId,
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "name", Value: 1}})
	opts.SetLimit(50)
	opts.SetSkip(0)

	ctx := context.Background()
	cur, err := nsSvc.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("failed to list namespace: %v", err)
		return nil, err
	}
	defer CloseCursor(ctx, cur)
	var nsList []*services.Namespace

	for cur.Next(ctx) {
		var ns services.Namespace
		err = cur.Decode(&ns)
		if err != nil {
			log.Printf("failed to decode while list namespace: %v", err)
			continue
		}
		nsList = append(nsList, &ns)
	}

	return nsList, nil
}

func (nsSvc *NamespaceService) UpdateNamespace(orgId string, nsId string, nsName string, schemaVersion string) error {
	ns, err := nsSvc.GetNamespaceById(orgId, nsId)
	if err != nil {
		log.Printf("failed to update namespace: %v", err)
		return err
	}

	filter := bson.M{
		"organizationid": orgId,
		"id":             nsId,
	}

	ns.Name = nsName
	ns.SchemaVersion = schemaVersion

	ctx := context.Background()
	result := nsSvc.collection.FindOneAndUpdate(ctx, filter, ns)
	if result.Err() != nil {
		log.Printf("failed to update namespace: %v", result.Err())
		return result.Err()
	}
	return result.Err()
}

func (nsSvc *NamespaceService) DeleteNamespace(orgId string, nsId string) error {
	filter := bson.M{
		"organizationid": orgId,
		"id":             nsId,
	}

	ctx := context.Background()
	result, err := nsSvc.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to delete namespace: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return services.ErrNamespaceNotFound
	}

	return nil
}

func (nsSvc *NamespaceService) ListNamespaceVars(orgId string, nsId string) ([]*services.NamespaceVar, error) {
	filter := bson.M{
		"orgid":       orgId,
		"namespaceid": nsId,
	}

	opts := options.Find()
	opts.SetSort(bson.D{primitive.E{Key: "id", Value: 1}})
	opts.SetLimit(50)
	opts.SetSkip(0)

	ctx := context.Background()
	cur, err := nsSvc.varCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Failed to list namespace vars: %v", err)
		return nil, err
	}

	defer CloseCursor(ctx, cur)

	var results []*services.NamespaceVar
	for cur.Next(ctx) {
		var nsVar services.NamespaceVar
		err := cur.Decode(&nsVar)
		if err != nil {
			log.Printf("Failed to decode namespace variable: %v", err)
			continue
		}
		results = append(results, &nsVar)
	}

	return results, nil
}

func (nsSvc *NamespaceService) GetVariablesAsMap(orgId string, nsId string) (map[string]string, error) {
	items, err := nsSvc.ListNamespaceVars(orgId, nsId)
	if err != nil {
		log.Printf("failed getting list of namespace variables for map: %v", err)
		return nil, err
	}

	result := make(map[string]string)
	for _, i := range items {
		result[i.Key] = i.Value
	}

	return result, nil
}

func (nsSvc *NamespaceService) GetNamespaceVariable(orgId string, nsId string, varId string) (*services.NamespaceVar, error) {
	filter := bson.M{
		"id":             varId,
		"organizationid": orgId,
		"namespaceid":    nsId,
	}

	ctx := context.Background()
	result := nsSvc.varCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		log.Printf("Failed to get namespace variable: %v", result.Err())
		return nil, result.Err()
	}

	var nsVar services.NamespaceVar
	err := result.Decode(&nsVar)
	if err != nil {
		log.Printf("Failed to get namespace variable: %v", err)
		return nil, result.Err()
	}

	return &nsVar, nil
}

func (nsSvc *NamespaceService) CreateNamespaceVariable(orgId string, nsId string, key string, value string) (*services.NamespaceVar, error) {
	newVar := &services.NamespaceVar{
		OrgId:       orgId,
		NamespaceId: nsId,
		Key:         key,
		Value:       value,
	}

	ctx := context.Background()
	_, err := nsSvc.varCollection.InsertOne(ctx, newVar)
	if err != nil {
		log.Printf("failed to create namespace varibale: %v", err)
		return nil, err
	}
	return newVar, nil
}

func (nsSvc *NamespaceService) UpdateNamespaceVariable(orgId string, nsId string, varId string, value string) (*services.NamespaceVar, error) {
	oldVar, err := nsSvc.GetNamespaceVariable(orgId, nsId, varId)
	if err != nil {
		log.Printf("failed to get namespace var for update: %v", err)
		return nil, err
	}

	oldVar.Value = value

	filter := bson.M{
		"organizationid": orgId,
		"namespaceid":    nsId,
		"key":            varId,
	}

	ctx := context.Background()
	result := nsSvc.varCollection.FindOneAndReplace(ctx, filter, oldVar)

	var nsVar services.NamespaceVar
	err = result.Decode(&nsVar)
	if err != nil {
		log.Printf("failed to update namespace variable: %v", err)
		return nil, result.Err()
	}

	return &nsVar, nil
}

func (nsSvc *NamespaceService) DeleteNamespaceVariable(orgId string, nsId string, varId string) error {
	filter := bson.M{
		"organizationid": orgId,
		"namespaceid":    nsId,
		"key":            varId,
	}

	ctx := context.Background()
	result := nsSvc.varCollection.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		log.Printf("failed to delete namespace variable: %v", result.Err())
		return result.Err()
	}

	return nil
}

func (nsSvc *NamespaceService) NamespaceVariableExists(orgId string, nsId string, varId string) bool {
	filter := bson.M{
		"organizationid": orgId,
		"namespaceid":    nsId,
		"key":            varId,
	}

	ctx := context.Background()
	result := nsSvc.varCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false
		}
		log.Printf("failed to find namespace variable: %v", result.Err())
		return false
	}

	return true
}
