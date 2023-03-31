/*
Copyright © 2023 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package common

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// used to return to the api client
type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// upload or form structs
type Dbentry interface {
	Getcontext(j []byte) (string, error)
	Marshal() ([]byte, error)
	MatchExpire(r *regexp.Regexp) bool
	MatchDescription(r *regexp.Regexp) bool
	MatchFile(r *regexp.Regexp) bool
	MatchCreated(r *regexp.Regexp) bool
	IsType(t int) bool
}

type Upload struct {
	Type        int       `json:"type"`
	Id          string    `json:"id"`
	Expire      string    `json:"expire"`
	File        string    `json:"file"`    // final filename (visible to the downloader)
	Members     []string  `json:"members"` // contains multiple files, so File is an archive
	Created     Timestamp `json:"uploaded"`
	Context     string    `json:"context"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
}

// this one is also used for marshalling to the client
type Response struct {
	Uploads []*Upload `json:"uploads"`
	Forms   []*Form   `json:"forms"`

	// integrate the Result struct so we can signal success
	Result
}

type Form struct {
	// Note the dual use of the Id: it will be used as onetime api key
	// from generated upload forms and  stored in the session store so
	// that the upload handler is able to check if the form object has
	// to be deleted immediately (if  its expire field has been set to
	// asap)
	Type        int       `json:"type"`
	Id          string    `json:"id"`
	Expire      string    `json:"expire"`
	Description string    `json:"description"`
	Created     Timestamp `json:"uploaded"`
	Context     string    `json:"context"`
	Url         string    `json:"url"`
	Notify      string    `json:"notify"`
}

const (
	TypeUpload = iota
	TypeForm
)

/*
   implement Dbentry interface
*/
func (upload Upload) Getcontext(j []byte) (string, error) {
	if err := json.Unmarshal(j, &upload); err != nil {
		return "", fmt.Errorf("unable to unmarshal json: %s", err)
	}

	return upload.Context, nil
}

func (form Form) Getcontext(j []byte) (string, error) {
	if err := json.Unmarshal(j, &form); err != nil {
		return "", fmt.Errorf("unable to unmarshal json: %s", err)
	}

	return form.Context, nil
}

func (upload Upload) Marshal() ([]byte, error) {
	jsonentry, err := json.Marshal(upload)
	if err != nil {
		return nil, fmt.Errorf("json marshalling failure: %s", err)
	}

	return jsonentry, nil
}

func (form Form) Marshal() ([]byte, error) {
	jsonentry, err := json.Marshal(form)
	if err != nil {
		return nil, fmt.Errorf("json marshalling failure: %s", err)
	}

	return jsonentry, nil
}

func (upload Upload) MatchExpire(r *regexp.Regexp) bool {
	return r.MatchString(upload.Expire)
}

func (upload Upload) MatchDescription(r *regexp.Regexp) bool {
	return r.MatchString(upload.Description)
}

func (upload Upload) MatchCreated(r *regexp.Regexp) bool {
	return r.MatchString(upload.Created.Time.String())
}

func (upload Upload) MatchFile(r *regexp.Regexp) bool {
	return r.MatchString(upload.File)
}

func (form Form) MatchExpire(r *regexp.Regexp) bool {
	return r.MatchString(form.Expire)
}

func (form Form) MatchDescription(r *regexp.Regexp) bool {
	return r.MatchString(form.Description)
}

func (form Form) MatchCreated(r *regexp.Regexp) bool {
	return r.MatchString(form.Created.Time.String())
}

func (form Form) MatchFile(r *regexp.Regexp) bool {
	return false
}

func (upload Upload) IsType(t int) bool {
	if upload.Type == t {
		return true
	}
	return false
}

func (form Form) IsType(t int) bool {
	if form.Type == t {
		return true
	}
	return false
}

/*
   Response methods
*/
func (r *Response) Append(entry Dbentry) {
	switch entry.(type) {
	case *Upload:
		r.Uploads = append(r.Uploads, entry.(*Upload))
	case Upload:
		r.Uploads = append(r.Uploads, entry.(*Upload))
	case Form:
		r.Forms = append(r.Forms, entry.(*Form))
	case *Form:
		r.Forms = append(r.Forms, entry.(*Form))
	default:
		panic("unknown type!")
	}
}

/*
   Extract  context, whatever  kind entry  is,  but we  don't know  in
   advance, only  after unmarshalling.  So try  Upload first,  if that
   fails, try Form.
*/
func GetContext(j []byte) (string, error) {
	upload := &Upload{}
	entryContext, err := upload.Getcontext(j)
	if err != nil {
		form := &Form{}
		entryContext, err = form.Getcontext(j)
		if err != nil {
			return "", fmt.Errorf("unable to unmarshal json: %s", err)
		}
	}

	return entryContext, nil
}

func Unmarshal(j []byte, t int) (Dbentry, error) {
	if t == TypeUpload {
		upload := &Upload{}
		if err := json.Unmarshal(j, &upload); err != nil {
			return upload, fmt.Errorf("unable to unmarshal json: %s", err)
		}
		return upload, nil
	} else {
		form := &Form{}
		if err := json.Unmarshal(j, &form); err != nil {
			return form, fmt.Errorf("unable to unmarshal json: %s", err)
		}
		return form, nil
	}
}
