module github.com/aurva-io/aurva/libraries/paginate

go 1.24

replace (
	github.com/aurva-io/aurva/audit => ../../audit
	github.com/aurva-io/aurva/avcontext => ../../avcontext
	github.com/aurva-io/aurva/avtime => ../../avtime
	github.com/aurva-io/aurva/command => ../../command
	github.com/aurva-io/aurva/commons => ../../commons
	github.com/aurva-io/aurva/company => ../../company
	github.com/aurva-io/aurva/dbresolver => ../../dbresolver
	github.com/aurva-io/aurva/dsorchestrator => ../../dsorchestrator
	github.com/aurva-io/aurva/elasticsearch/client => ../../elasticsearch/client
	github.com/aurva-io/aurva/elasticsearch/reporting => ../../elasticsearch/reporting
	github.com/aurva-io/aurva/healthchecker => ../../healthchecker
	github.com/aurva-io/aurva/logging => ../../logging
)
