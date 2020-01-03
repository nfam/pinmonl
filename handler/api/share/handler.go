package share

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func HandleList(shares store.ShareStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ml, err := shares.List(r.Context(), &store.ShareOpts{})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := make([]interface{}, len(ml))
		for i, m := range ml {
			resp[i] = Resp(m)
		}
		response.JSON(w, resp)
	}
}

func HandleFind(shares store.ShareStore, sharetags store.ShareTagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := request.ShareFrom(ctx)

		mtsm, err := sharetags.ListTags(ctx, &store.ShareTagOpts{ShareID: m.ID, Kind: model.MustTag})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		mts := mtsm[m.ID]

		atsm, err := sharetags.ListTags(ctx, &store.ShareTagOpts{ShareID: m.ID, Kind: model.AnyTag})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		ats := atsm[m.ID]

		response.JSON(w, DetailResp(m, mts, ats))
	}
}

func HandleCreate(shares store.ShareStore, sharetags store.ShareTagStore, tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.Validate()
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		var m model.Share
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = shares.Create(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mts, err := tags.List(ctx, &store.TagOpts{Names: in.MustTags})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = sharetags.ReAssocTags(ctx, m, model.MustTag, mts)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ats := make([]model.Tag, 0)
		if len(in.AnyTags) > 0 {
			ats, err = tags.List(ctx, &store.TagOpts{Names: in.AnyTags})
			if err != nil {
				response.InternalError(w, err)
				return
			}
			err = sharetags.ReAssocTags(ctx, m, model.AnyTag, ats)
			if err != nil {
				response.InternalError(w, err)
				return
			}
		}

		response.JSON(w, DetailResp(m, mts, ats))
	}
}

func HandleUpdate(shares store.ShareStore, sharetags store.ShareTagStore, tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.Validate()
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		m, _ := request.ShareFrom(ctx)
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = shares.Create(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mts, err := tags.List(ctx, &store.TagOpts{Names: in.MustTags})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = sharetags.ReAssocTags(ctx, m, model.MustTag, mts)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ats := make([]model.Tag, 0)
		if len(in.AnyTags) > 0 {
			ats, err = tags.List(ctx, &store.TagOpts{Names: in.AnyTags})
			if err != nil {
				response.InternalError(w, err)
				return
			}
			err = sharetags.ReAssocTags(ctx, m, model.AnyTag, ats)
			if err != nil {
				response.InternalError(w, err)
				return
			}
		}

		response.JSON(w, DetailResp(m, mts, ats))
	}
}

func HandleDelete(shares store.ShareStore, sharetags store.ShareTagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := request.ShareFrom(ctx)

		err := shares.Delete(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = sharetags.ClearByKind(ctx, m, model.MustTag)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = sharetags.ClearByKind(ctx, m, model.AnyTag)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.NoContent(w)
	}
}