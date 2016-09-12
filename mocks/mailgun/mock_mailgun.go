package mockmailgun

import (
	"io"
	"net/http"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type MockMailGun struct{}

func (mmg *MockMailGun) Domain() string {
	panic("Not implemented")
}
func (mmg *MockMailGun) ApiKey() string {
	panic("Not implemented")
}
func (mmg *MockMailGun) PublicApiKey() string {
	panic("Not implemented")
}
func (mmg *MockMailGun) Client() *http.Client {
	panic("Not implemented")
}
func (mmg *MockMailGun) SetClient(client *http.Client) {
	panic("Not implemented")
}
func (mmg *MockMailGun) Send(m *mailgun.Message) (string, string, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) ValidateEmail(email string) (mailgun.EmailVerification, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) ParseAddresses(addresses ...string) ([]string, []string, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetBounces(limit, skip int) (int, []mailgun.Bounce, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetSingleBounce(address string) (mailgun.Bounce, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) AddBounce(address, code, error string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteBounce(address string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetStats(limit int, skip int, startDate *time.Time, event ...string) (int, []mailgun.Stat, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteTag(tag string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetDomains(limit, skip int) (int, []mailgun.Domain, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetSingleDomain(domain string) (mailgun.Domain, []mailgun.DNSRecord, []mailgun.DNSRecord, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateDomain(name string, smtpPassword string, spamAction string, wildcard bool) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteDomain(name string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetCampaigns() (int, []mailgun.Campaign, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateCampaign(name, id string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) UpdateCampaign(oldId, name, newId string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteCampaign(id string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetComplaints(limit, skip int) (int, []mailgun.Complaint, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetSingleComplaint(address string) (mailgun.Complaint, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetStoredMessage(id string) (mailgun.StoredMessage, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetStoredMessageRaw(id string) (mailgun.StoredMessageRaw, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteStoredMessage(id string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetCredentials(limit, skip int) (int, []mailgun.Credential, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateCredential(login, password string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) ChangeCredentialPassword(id, password string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteCredential(id string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetUnsubscribes(limit, skip int) (int, []mailgun.Unsubscription, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetUnsubscribesByAddress(string) (int, []mailgun.Unsubscription, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) Unsubscribe(address, tag string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) RemoveUnsubscribe(string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) RemoveUnsubscribeWithTag(a, t string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateComplaint(string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteComplaint(string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetRoutes(limit, skip int) (int, []mailgun.Route, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetRouteByID(string) (mailgun.Route, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateRoute(mailgun.Route) (mailgun.Route, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteRoute(string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) UpdateRoute(string, mailgun.Route) (mailgun.Route, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetWebhooks() (map[string]string, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateWebhook(kind, url string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteWebhook(kind string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetWebhookByType(kind string) (string, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) UpdateWebhook(kind, url string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetLists(limit, skip int, filter string) (int, []mailgun.List, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateList(mailgun.List) (mailgun.List, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteList(string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetListByAddress(string) (mailgun.List, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) UpdateList(string, mailgun.List) (mailgun.List, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetMembers(limit, skip int, subfilter *bool, address string) (int, []mailgun.Member, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) GetMemberByAddress(MemberAddr, listAddr string) (mailgun.Member, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateMember(merge bool, addr string, prototype mailgun.Member) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) CreateMemberList(subscribed *bool, addr string, newMembers []interface{}) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) UpdateMember(Member, list string, prototype mailgun.Member) (mailgun.Member, error) {
	panic("Not implemented")
}
func (mmg *MockMailGun) DeleteMember(Member, list string) error {
	panic("Not implemented")
}
func (mmg *MockMailGun) NewMessage(from, subject, text string, to ...string) *mailgun.Message {
	panic("Not implemented")
}
func (mmg *MockMailGun) NewMIMEMessage(body io.ReadCloser, to ...string) *mailgun.Message {
	panic("Not implemented")
}
func (mmg *MockMailGun) NewEventIterator() *mailgun.EventIterator {
	panic("Not implemented")
}
