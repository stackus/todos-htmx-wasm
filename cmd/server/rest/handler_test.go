package rest

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/mock"

	"github.com/stackus/todos-htmx-wasm/internal/domain"
	"github.com/stackus/todos-htmx-wasm/internal/features/todos"
)

func Test_handler_Create(t *testing.T) {
	var todo = &domain.Todo{
		ID:          uuid.New(),
		Description: "first",
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	type fields struct {
		todosSvc *todos.MockService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		args           args
		mock           func(f fields)
		wantStatusCode int
		wantHeader     http.Header
		want           *domain.Todo
	}{
		"Create": {
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"description":"first"}`))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
			},
			mock: func(f fields) {
				f.todosSvc.EXPECT().Add(context.Background(), "first").Return(todo, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantHeader:     http.Header{},
			want:           todo,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				todosSvc: todos.NewMockService(t),
			}
			h := handler{
				todosSvc: f.todosSvc,
			}
			if tt.mock != nil {
				tt.mock(f)
			}

			h.Create(tt.args.w, tt.args.r)

			res := tt.args.w.(*httptest.ResponseRecorder)
			if res.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("handler.Create() StatusCode = %v, want %v", res.Result().StatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(res.Result().Header, tt.wantHeader) {
				t.Errorf("handler.Create() Header = %v, want %v", res.Result().Header, tt.wantHeader)
			}

			gotJSON := res.Body.Bytes()
			if tt.want == nil {
				if len(gotJSON) > 0 {
					t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), "")
				}
				return
			}
			var gotTodo domain.Todo
			if err := json.Unmarshal(gotJSON, &gotTodo); err != nil {
				t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), tt.want)
			}
			if gotTodo.ID != tt.want.ID || gotTodo.Description != tt.want.Description || gotTodo.Completed != tt.want.Completed {
				t.Errorf("handler.Create() Body = %v, want %v", gotTodo, tt.want)
			}
		})
	}
}

func Test_handler_Delete(t *testing.T) {
	var todoID = uuid.New()
	type fields struct {
		todosSvc *todos.MockService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		args           args
		mock           func(f fields)
		wantStatusCode int
		wantHeader     http.Header
		want           *domain.Todo
	}{
		"Delete": {
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodDelete, "/", nil)
					rCtx := chi.NewRouteContext()
					rCtx.URLParams.Add("todoId", todoID.String())
					req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
					return req
				}(),
			},
			mock: func(f fields) {
				f.todosSvc.EXPECT().Remove(mock.AnythingOfType("*context.valueCtx"), todoID).Return(nil)
			},
			wantStatusCode: http.StatusNoContent,
			wantHeader:     http.Header{},
			want:           nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				todosSvc: todos.NewMockService(t),
			}
			h := handler{
				todosSvc: f.todosSvc,
			}
			if tt.mock != nil {
				tt.mock(f)
			}

			h.Delete(tt.args.w, tt.args.r)

			res := tt.args.w.(*httptest.ResponseRecorder)
			if res.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("handler.Delete() StatusCode = %v, want %v", res.Result().StatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(res.Result().Header, tt.wantHeader) {
				t.Errorf("handler.Delete() Header = %v, want %v", res.Result().Header, tt.wantHeader)
			}
			gotJSON := res.Body.Bytes()
			if tt.want == nil {
				if len(gotJSON) > 0 {
					t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), "")
				}
				return
			}
			var gotTodo domain.Todo
			if err := json.Unmarshal(gotJSON, &gotTodo); err != nil {
				t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), tt.want)
			}
			if gotTodo.ID != tt.want.ID || gotTodo.Description != tt.want.Description || gotTodo.Completed != tt.want.Completed {
				t.Errorf("handler.Delete() Body = %v, want %v", gotTodo, tt.want)
			}
		})
	}
}

func Test_handler_Get(t *testing.T) {
	var todoID = uuid.New()
	var todo = &domain.Todo{
		ID:          todoID,
		Description: "test",
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	type fields struct {
		todosSvc *todos.MockService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		args           args
		mock           func(f fields)
		wantStatusCode int
		wantHeader     http.Header
		want           *domain.Todo
	}{
		"Get": {
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/", nil)
					rCtx := chi.NewRouteContext()
					rCtx.URLParams.Add("todoId", todoID.String())
					req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
					return req
				}(),
			},
			mock: func(f fields) {
				f.todosSvc.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), todoID).Return(todo, nil)
			},
			wantStatusCode: http.StatusOK,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			want: todo,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				todosSvc: todos.NewMockService(t),
			}
			h := handler{
				todosSvc: f.todosSvc,
			}
			if tt.mock != nil {
				tt.mock(f)
			}

			h.Get(tt.args.w, tt.args.r)

			res := tt.args.w.(*httptest.ResponseRecorder)
			if res.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("handler.Get() StatusCode = %v, want %v", res.Result().StatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(res.Result().Header, tt.wantHeader) {
				t.Errorf("handler.Get() Header = %v, want %v", res.Result().Header, tt.wantHeader)
			}
			gotJSON := res.Body.Bytes()
			if tt.want == nil {
				if len(gotJSON) > 0 {
					t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), "")
				}
				return
			}
			var gotTodo domain.Todo
			if err := json.Unmarshal(gotJSON, &gotTodo); err != nil {
				t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), tt.want)
			}
			if gotTodo.ID != tt.want.ID || gotTodo.Description != tt.want.Description || gotTodo.Completed != tt.want.Completed {
				t.Errorf("handler.Get() Body = %v, want %v", gotTodo, tt.want)
			}
		})
	}
}

func Test_handler_Search(t *testing.T) {
	var todo = &domain.Todo{
		ID:          uuid.New(),
		Description: "test",
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	type fields struct {
		todosSvc *todos.MockService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		args           args
		mock           func(f fields)
		wantStatusCode int
		wantHeader     http.Header
		want           []*domain.Todo
	}{
		"Search": {
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/?search=test", nil)
					return req
				}(),
			},
			mock: func(f fields) {
				f.todosSvc.EXPECT().Search(mock.AnythingOfType("*context.emptyCtx"), "test").Return([]*domain.Todo{todo}, nil)
			},
			wantStatusCode: http.StatusOK,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			want: []*domain.Todo{todo},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				todosSvc: todos.NewMockService(t),
			}
			h := handler{
				todosSvc: f.todosSvc,
			}
			if tt.mock != nil {
				tt.mock(f)
			}

			h.Search(tt.args.w, tt.args.r)

			res := tt.args.w.(*httptest.ResponseRecorder)
			if res.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("handler.Search() StatusCode = %v, want %v", res.Result().StatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(res.Result().Header, tt.wantHeader) {
				t.Errorf("handler.Search() Header = %v, want %v", res.Result().Header, tt.wantHeader)
			}
			gotJSON := res.Body.Bytes()
			if tt.want == nil {
				if len(gotJSON) > 0 {
					t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), "")
				}
				return
			}
			var gotTodos []*domain.Todo
			if err := json.Unmarshal(gotJSON, &gotTodos); err != nil {
				t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), tt.want)
			}
			if len(gotTodos) != len(tt.want) {
				t.Errorf("handler.Search() Body = %v, want %v", gotTodos, tt.want)
			}
			for i, gotTodo := range gotTodos {
				if gotTodo.ID != tt.want[i].ID || gotTodo.Description != tt.want[i].Description || gotTodo.Completed != tt.want[i].Completed {
					t.Errorf("handler.Search() Body = %v, want %v", gotTodo, tt.want)
				}
			}
		})
	}
}

func Test_handler_Sort(t *testing.T) {
	var firstTodo = &domain.Todo{
		ID:          uuid.New(),
		Description: "test",
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	var secondTodo = &domain.Todo{
		ID:          uuid.New(),
		Description: "test2",
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	var todoIDs = []uuid.UUID{firstTodo.ID, secondTodo.ID}
	var ids = struct {
		IDs []uuid.UUID `json:"ids"`
	}{
		IDs: todoIDs,
	}
	data, _ := json.Marshal(ids)
	type fields struct {
		todosSvc *todos.MockService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		args           args
		mock           func(f fields)
		wantStatusCode int
		wantHeader     http.Header
		want           []*domain.Todo
	}{
		"Sort": {
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					ids := make([]string, len(todoIDs))
					for i, id := range todoIDs {
						ids[i] = id.String()
					}
					req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
					req.Header.Set("Content-Type", "application/json")
					return req
				}(),
			},
			mock: func(f fields) {
				f.todosSvc.EXPECT().Sort(mock.AnythingOfType("*context.emptyCtx"), todoIDs).Return([]*domain.Todo{
					firstTodo, secondTodo,
				}, nil)
			},
			wantStatusCode: http.StatusOK,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			want: []*domain.Todo{firstTodo, secondTodo},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				todosSvc: todos.NewMockService(t),
			}
			h := handler{
				todosSvc: f.todosSvc,
			}
			if tt.mock != nil {
				tt.mock(f)
			}

			h.Sort(tt.args.w, tt.args.r)

			res := tt.args.w.(*httptest.ResponseRecorder)
			if res.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("handler.Sort() StatusCode = %v, want %v", res.Result().StatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(res.Result().Header, tt.wantHeader) {
				t.Errorf("handler.Sort() Header = %v, want %v", res.Result().Header, tt.wantHeader)
			}
			gotJSON := res.Body.Bytes()
			if tt.want == nil {
				if len(gotJSON) > 0 {
					t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), "")
				}
				return
			}
			var gotTodos []*domain.Todo
			if err := json.Unmarshal(gotJSON, &gotTodos); err != nil {
				t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), tt.want)
			}
			if len(gotTodos) != len(tt.want) {
				t.Errorf("handler.Search() Body = %v, want %v", gotTodos, tt.want)
			}
			for i, gotTodo := range gotTodos {
				if gotTodo.ID != tt.want[i].ID || gotTodo.Description != tt.want[i].Description || gotTodo.Completed != tt.want[i].Completed {
					t.Errorf("handler.Search() Body = %v, want %v", gotTodo, tt.want)
				}
			}
		})
	}
}

func Test_handler_Update(t *testing.T) {
	var todoID = uuid.New()
	var updated = &domain.Todo{
		ID:          todoID,
		Description: "first",
		Completed:   true,
		CreatedAt:   time.Now(),
	}
	data, _ := json.Marshal(updated)
	type fields struct {
		todosSvc *todos.MockService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		args           args
		mock           func(f fields)
		wantStatusCode int
		wantHeader     http.Header
		want           *domain.Todo
	}{
		"Update": {
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
					req.Header.Set("Content-Type", "application/json")
					rCtx := chi.NewRouteContext()
					rCtx.URLParams.Add("todoId", todoID.String())
					return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))
				}(),
			},
			mock: func(f fields) {
				f.todosSvc.EXPECT().Update(mock.AnythingOfType("*context.valueCtx"), todoID, true, "first").Return(updated, nil)
			},
			wantStatusCode: http.StatusOK,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			want: updated,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := fields{
				todosSvc: todos.NewMockService(t),
			}
			h := handler{
				todosSvc: f.todosSvc,
			}
			if tt.mock != nil {
				tt.mock(f)
			}

			h.Update(tt.args.w, tt.args.r)

			res := tt.args.w.(*httptest.ResponseRecorder)
			if res.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("handler.Update() StatusCode = %v, want %v", res.Result().StatusCode, tt.wantStatusCode)
			}
			if !reflect.DeepEqual(res.Result().Header, tt.wantHeader) {
				t.Errorf("handler.Update() Header = %v, want %v", res.Result().Header, tt.wantHeader)
			}
			gotJSON := res.Body.Bytes()
			if tt.want == nil {
				if len(gotJSON) > 0 {
					t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), "")
				}
				return
			}
			var gotTodo domain.Todo
			if err := json.Unmarshal(gotJSON, &gotTodo); err != nil {
				t.Errorf("handler.Create() Body = %v, want %v", string(gotJSON), tt.want)
			}
			if gotTodo.ID != tt.want.ID || gotTodo.Description != tt.want.Description || gotTodo.Completed != tt.want.Completed {
				t.Errorf("handler.Get() Body = %v, want %v", gotTodo, tt.want)
			}
		})
	}
}
