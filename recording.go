package twilio

import (
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

type RecordingService struct {
	client *Client
}

const recordingsPathPart = "Recordings"

type Recording struct {
	Sid         string         `json:"sid"`
	Duration    TwilioDuration `json:"duration"`
	CallSid     string         `json:"call_sid"`
	Status      Status         `json:"status"`
	Price       string         `json:"price"`
	PriceUnit   string         `json:"price_unit"`
	DateCreated TwilioTime     `json:"date_created"`
	AccountSid  string         `json:"account_sid"`
	APIVersion  string         `json:"api_version"`
	Channels    uint           `json:"channels"`
	DateUpdated TwilioTime     `json:"date_updated"`
	URI         string         `json:"uri"`
}

// URL returns the URL that can be used to play this recording, based on the
// extension. No error is returned if you provide an invalid extension. As of
// October 2016, the valid values are ".wav" and ".mp3".
func (r *Recording) URL(extension string) string {
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	return strings.Join([]string{BaseURL, r.APIVersion, "Accounts", r.AccountSid, recordingsPathPart, r.Sid + extension}, "/")
}

// FriendlyPrice flips the sign of the Price (which is usually reported from
// the API as a negative number) and adds an appropriate currency symbol in
// front of it. For example, a PriceUnit of "USD" and a Price of "-1.25" is
// reported as "$1.25".
func (r *Recording) FriendlyPrice() string {
	if r == nil {
		return ""
	}
	return price(r.PriceUnit, r.Price)
}

type RecordingPage struct {
	Page
	Recordings []*Recording
}

func (r *RecordingService) Get(ctx context.Context, sid string) (*Recording, error) {
	recording := new(Recording)
	err := r.client.GetResource(ctx, recordingsPathPart, sid, recording)
	return recording, err
}

func (r *RecordingService) GetPage(ctx context.Context, data url.Values) (*RecordingPage, error) {
	rp := new(RecordingPage)
	err := r.client.ListResource(ctx, recordingsPathPart, data, rp)
	return rp, err
}

type RecordingPageIterator struct {
	p *PageIterator
}

// GetPageIterator returns an iterator which can be used to retrieve pages.
func (r *RecordingService) GetPageIterator(data url.Values) *RecordingPageIterator {
	iter := NewPageIterator(r.client, data, recordingsPathPart)
	return &RecordingPageIterator{
		p: iter,
	}
}

// Next returns the next page of resources. If there are no more resources,
// NoMoreResults is returned.
func (r *RecordingPageIterator) Next(ctx context.Context) (*RecordingPage, error) {
	rp := new(RecordingPage)
	err := r.p.Next(ctx, rp)
	if err != nil {
		return nil, err
	}
	r.p.SetNextPageURI(rp.NextPageURI)
	return rp, nil
}