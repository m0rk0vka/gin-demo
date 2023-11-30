// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/m0rk0vka/gin-demo.
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Vladimir Tatosyan<v.tatosyan@yandex.ru>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run()
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
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
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	 '200':
//		description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
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
//	   '404':
//			  description: invalid recipe id
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i, r := range recipes {
		if r.ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipe.ID = id
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
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
//	   '404':
//			  description: invalid recipe id
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	index := -1
	for i, r := range recipes {
		if r.ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
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
func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		found := false
		for _, t := range recipe.Tags {
			if t == tag {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipe)

		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
