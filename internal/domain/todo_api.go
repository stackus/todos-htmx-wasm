package domain

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type TodoApi struct {
	client *http.Client
	host   string
}

var _ TodoRepository = (*TodoApi)(nil)

func NewTodoApi(host string) *TodoApi {
	return &TodoApi{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		host: host,
	}
}

func (t *TodoApi) Add(description string) (*Todo, error) {
	type addTodoRequest struct {
		Description string `json:"description"`
	}

	data, err := json.Marshal(addTodoRequest{Description: description})
	if err != nil {
		return nil, ErrMarshaling{Err: err}
	}

	resp, err := t.doRequest(http.MethodPost, "/todos", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	type addTodoResponse *Todo

	var addTodoResp addTodoResponse
	err = json.NewDecoder(resp.Body).Decode(&addTodoResp)
	if err != nil {
		return nil, ErrUnmarshaling{Err: err}
	}

	return addTodoResp, nil
}

func (t *TodoApi) Remove(id uuid.UUID) error {
	_, err := t.doRequest(http.MethodDelete, "/todos/"+id.String(), nil)
	if err != nil {
		return err
	}

	return nil
}

func (t *TodoApi) Update(id uuid.UUID, completed bool, description string) (*Todo, error) {
	type updateTodoRequest struct {
		Completed   bool   `json:"completed"`
		Description string `json:"description"`
	}

	data, err := json.Marshal(updateTodoRequest{Completed: completed, Description: description})
	if err != nil {
		return nil, ErrMarshaling{Err: err}
	}

	resp, err := t.doRequest(http.MethodPatch, "/todos/"+id.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	type updateTodoResponse *Todo

	var updateTodoResp updateTodoResponse
	err = json.NewDecoder(resp.Body).Decode(&updateTodoResp)
	if err != nil {
		return nil, ErrUnmarshaling{Err: err}
	}

	return updateTodoResp, nil
}

func (t *TodoApi) Search(search string) ([]*Todo, error) {
	resp, err := t.doRequest(http.MethodGet, "/todos", nil)
	if err != nil {
		return nil, err
	}

	type searchTodoResponse []*Todo

	var searchTodoResp searchTodoResponse
	err = json.NewDecoder(resp.Body).Decode(&searchTodoResp)
	if err != nil {
		return nil, ErrUnmarshaling{Err: err}
	}

	return searchTodoResp, nil
}

func (t *TodoApi) All() ([]*Todo, error) {
	return t.Search("")
}

func (t *TodoApi) Get(id uuid.UUID) (*Todo, error) {
	resp, err := t.doRequest(http.MethodGet, "/todos/"+id.String(), nil)
	if err != nil {
		return nil, err
	}

	type getTodoResponse *Todo

	var getTodoResp getTodoResponse
	err = json.NewDecoder(resp.Body).Decode(&getTodoResp)
	if err != nil {
		return nil, ErrUnmarshaling{Err: err}
	}

	return getTodoResp, nil
}

func (t *TodoApi) Reorder(ids []uuid.UUID) ([]*Todo, error) {
	type reorderTodoRequest struct {
		IDs []uuid.UUID `json:"ids"`
	}

	data, err := json.Marshal(reorderTodoRequest{IDs: ids})
	if err != nil {
		return nil, ErrMarshaling{Err: err}
	}

	resp, err := t.doRequest(http.MethodPost, "/todos/sort", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	type reorderTodoResponse []*Todo

	var reorderTodoResp reorderTodoResponse
	err = json.NewDecoder(resp.Body).Decode(&reorderTodoResp)
	if err != nil {
		return nil, ErrUnmarshaling{Err: err}
	}

	return reorderTodoResp, nil
}

func (t *TodoApi) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, t.host+path, body)
	if err != nil {
		return nil, ErrCreateRequest{Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, ErrMakeRequest{Err: err}
	}
	return resp, nil
}
