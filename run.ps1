echo .\server.go .\client.go .\rate-limiter.go | 
  % { $_ | start-threadjob { go run $input } } | 
  receive-job -wait -auto | ft -a
