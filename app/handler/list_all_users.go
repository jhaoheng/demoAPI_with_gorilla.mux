package handler

import (
	"app/models"
	"app/modules"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

type ListAllUsers struct {
	query                *ListAllUsersQuery
	path                 *ListAllUsersPath
	body                 *ListAllUsersBody
	model_get_all_counts models.IUser
	model_get_users      models.IUser
}

type ListAllUsersQuery struct {
	Paging  string
	Sorting string
}
type ListAllUsersPath struct{}
type ListAllUsersBody struct{}
type ListAllUsersResp struct {
	Total int                    `json:"total"`
	Users []ListAllUsersRespUser `json:"users"`
}

type ListAllUsersRespUser struct {
	Account   string `json:"account"`
	Fullname  string `json:"fullname"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ListAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	api := ListAllUsers{
		query: &ListAllUsersQuery{
			Paging:  r.URL.Query().Get("paging"),
			Sorting: r.URL.Query().Get("sorting"),
		},
		model_get_all_counts: models.NewUser(),
		model_get_users:      models.NewUser(),
	}
	resp, status, err := api.do()
	modules.NewResp(w, r).Set(modules.RespContect{Data: resp, Error: err, Stutus: status})
}

func (api *ListAllUsers) do() (*ListAllUsersResp, int, error) {

	if err := api.check_paging(); err != nil {
		return nil, http.StatusBadRequest, err
	}
	//
	if err := api.check_sorting(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	total, err := api.model_get_all_counts.GetAllCount()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	results, err := api.model_get_users.ListBy(api.query.Paging, api.query.Sorting, 10)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//
	datas := []ListAllUsersRespUser{}
	for _, user := range results {
		data := ListAllUsersRespUser{
			Account:   user.Acct,
			Fullname:  user.Fullname,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		}
		datas = append(datas, data)
	}

	resp := &ListAllUsersResp{
		Total: int(total),
		Users: datas,
	}
	return resp, http.StatusOK, nil
}

func (api *ListAllUsers) check_paging() error {
	if len(api.query.Paging) == 0 || strings.EqualFold(api.query.Paging, "0") {
		api.query.Paging = "1"
		return nil
	}
	var re = regexp.MustCompile(`[0-9]$`)
	if !re.MatchString(api.query.Paging) {
		return errors.New("paging must be number")
	}
	return nil
}

func (api *ListAllUsers) check_sorting() error {
	if len(api.query.Sorting) == 0 {
		api.query.Sorting = "asc"
	}

	var re = regexp.MustCompile(`^(asc|desc)$`)
	if !re.MatchString(api.query.Sorting) {
		return errors.New(`sorting must be 'asc' or 'desc'`)
	}
	return nil
}
