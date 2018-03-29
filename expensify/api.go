package expensify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocarina/gocsv"

	resty "gopkg.in/resty.v1"
)

const endpoint = "https://integrations.expensify.com/Integration-Server/ExpensifyIntegrations"
const exportExpensesTemplate = `
<#if addHeader == true>
    merchant,amount,created<#lt>
</#if>
<#list reports as report>
	<#list report.transactionList as expense>
		<#if expense.modifiedMerchant?has_content>
			<#assign merchant = expense.modifiedMerchant>
		<#else>
			<#assign merchant = expense.merchant>
		</#if>
		<#if expense.convertedAmount?has_content>
			<#assign amount = expense.convertedAmount>
		<#elseif expense.modifiedAmount?has_content>
			<#assign amount = expense.modifiedAmount>
		<#else>
			<#assign amount = expense.amount>
		</#if>
		<#if expense.modifiedCreated?has_content>
			<#assign created = expense.modifiedCreated>
		<#else>
			<#assign created = expense.created>
		</#if>
		${merchant},<#t>
		${amount},<#t>
		${created}<#lt>
    </#list>
</#list>`

var validCSVFile = regexp.MustCompile(`.*\.csv$`)

// API provides an interface to interact with the Expensify API
type API struct {
	partnerUserID     string
	partnerUserSecret string
}

// requestJobDescription is teh data carrier for each API request
type requestJobDescription struct {
	Type           string             `json:"type"`
	Credentials    requestCredentials `json:"credentials"`
	InputSettings  interface{}        `json:"inputSettings"`
	OutputSettings interface{}        `json:"outputSettings"`
	OnReceive      interface{}        `json:"onReceive"`
	FileName       string             `json:"fileName,omit"`
}

// requestCredentials is used for sending the credentials in a request
type requestCredentials struct {
	PartnerUserID     string `json:"partnerUserID"`
	PartnerUserSecret string `json:"partnerUserSecret"`
}

type Expense struct {
	Merchant string `csv:"merchant"`
	Amount   int    `csv:"amount"`
	Created  string `csv:"created"`
}

func New(id, secret string) API {
	return API{
		partnerUserID:     id,
		partnerUserSecret: secret,
	}
}

func (a *API) doCall(jobDesc requestJobDescription) (*resty.Response, error) {
	jobDesc.Credentials = requestCredentials{
		PartnerUserID:     a.partnerUserID,
		PartnerUserSecret: a.partnerUserSecret,
	}

	body, err := json.Marshal(jobDesc)
	if err != nil {
		return nil, err
	}
	return resty.SetHostURL(endpoint).R().SetFormData(map[string]string{"requestJobDescription": string(body), "template": exportExpensesTemplate}).Post(endpoint)
}

func (a *API) download(file string) ([]byte, error) {
	jobDesc := requestJobDescription{
		Type:     "download",
		FileName: file,
	}

	resp, err := a.doCall(jobDesc)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// GetExpenses gets all expenses from reports with a specific state, from a start and end date ans will respect a limit of reports checked
func (a *API) GetExpenses(reportState, start, end string, limit int64) ([]Expense, error) {
	jobDesc := requestJobDescription{
		Type: "file",
		OnReceive: map[string][]string{
			"immediateResponse": []string{"returnRandomFileName"},
		},
		InputSettings: map[string]interface{}{
			"type":        "combinedReportData",
			"reportState": reportState,
			"limit":       strconv.FormatInt(limit, 10),
			"filters": map[string]string{
				"startDate": start,
				"endDate":   end,
			},
		},
		OutputSettings: map[string]string{
			"fileExtension": "csv",
		},
	}

	rsp, err := a.doCall(jobDesc)
	if err != nil {
		return nil, err
	}

	fileName := string(rsp.Body())
	if !validCSVFile.MatchString(fileName) {
		return nil, fmt.Errorf("Expensify returned an error: %s", fileName)
	}
	file, err := a.download(fileName)
	if err != nil {
		return nil, err
	}

	expenses := []Expense{}

	if err := gocsv.Unmarshal(bytes.NewReader(file), &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}
