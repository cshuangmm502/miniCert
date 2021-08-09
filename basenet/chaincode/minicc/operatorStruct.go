package main

type Operator struct {
	SerialNumber string
	IssueWitnessIndex string
	WitnessIndex	int
	CertificateHash	string
	Histories	[]HistoryOperator
}

type HistoryOperator struct {
	TxId	string
	Operator	Operator
}