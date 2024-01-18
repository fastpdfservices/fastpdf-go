#
# Run fastpdf-go SDK tests
#
#

# Clean cache
go clean -testcache

# Run TEsts
go test ./fastpdf
