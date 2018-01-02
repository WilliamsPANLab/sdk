#include <Rcpp.h>
#include <flywheel.h>
using namespace Rcpp;

// Currently unused

// This file wraps the Flywheel C bridge in function signatures that R can digest.
// Ref: https://www.r-bloggers.com/three-ways-to-call-cc-from-r/

// [[Rcpp::export]]
int double_me3(int x) {
  // takes a numeric input and doubles it
  return 2 * x;
}

{{range .Signatures}}
// [[Rcpp::export]]
char* fw_{{.Name}}(char* apiKey, int root{{range .Params}}, {{.CType}} {{.Name}}{{end}}, int* status) {
	return {{.Name}}(apiKey{{range .Params}}, {{.Name}}{{end}}, status);
}
{{end}}
