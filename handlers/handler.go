package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m0rk0vka/gin-demo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	 '200':
//			description: Successful operation
//
//	 '500':
//
//		description: Error while executing find comand to collection
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cur.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

// swagger:operation POST /recipes recipes createRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//
//	   '200':
//	     description: successful operation
//	   '400':
//			  description: invalid input
func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	if _, err := handler.collection.InsertOne(handler.ctx, recipe); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting a new recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: id of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	   '200':
//	     description: successful operation
//	   '400':
//			  description: invalid input
//	   '500':
//			  description: error while update db
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{
		"name":         recipe.Name,
		"instructions": recipe.Instructions,
		"ingredients":  recipe.Ingredients,
		"tags":         recipe.Tags,
	}}
	opts := options.Update().SetUpsert(true)
	_, err := handler.collection.UpdateOne(handler.ctx, filter, update, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meddage": "Recipe has been updated"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: id of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	   '200':
//	     description: successful operation
//	   '500':
//			  description: error while delete recipe
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	if _, err := handler.collection.DeleteOne(handler.ctx, filter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes findRecipe
// Search reipes based on tag
// ---
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	 '200':
//		  description: successful operation
//	 '500':
//			  description: error while executing find command on collection
func (handler *RecipesHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]models.Recipe, 0)

	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	defer cur.Close(handler.ctx)

	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		find := false
		for _, t := range recipe.Tags {
			if t == tag {
				find = true
			}
		}
		if find {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
