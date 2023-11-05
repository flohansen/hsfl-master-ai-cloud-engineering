package recipes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/stockhut/hsfl-master-ai-cloud-engineering/authentication/middleware"
	"github.com/stockhut/hsfl-master-ai-cloud-engineering/common/htmx"
	"github.com/stockhut/hsfl-master-ai-cloud-engineering/common/presenter/json_presenter"
)

func (ctrl *Controller) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestBody createRecipeRequestBody
	if err := json.Unmarshal(body, &requestBody); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.JwtContextKey).(jwt.MapClaims)

	username, ok := claims["name"]
	if !ok {
		fmt.Println("failed to read name from jwt")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe := recipeRequestToModel(requestBody, username.(string))

	newRecipe, err := ctrl.repo.CreateRecipe(recipe)
	if err != nil {
		fmt.Printf("Failed to save recipe: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := recipeToResponseModel(newRecipe)

	if htmx.IsHtmxRequest(r) {
		tmplFile := "templates/CreateRecipe/recipeSuccessfulCreate.html"
		tmpl, err := template.ParseFiles(tmplFile)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		err = tmpl.Execute(w, response)
		if err != nil {
			panic(err)
		}
	} else {
		json_presenter.JsonPresenter(w, http.StatusOK, response)
	}

}
