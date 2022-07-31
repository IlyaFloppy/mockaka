package mockaka

/*
func NewJSONStub[I, O protomsg](t *testing.T, method string, request, response json.RawMessage) Stub {
	var (
		req I
		res O
	)

	if err := json.Unmarshal(request, req); err != nil {
		t.Fatalf("failed to unmarshal request: %s", err.Error())
	}

	if err := json.Unmarshal(response, res); err != nil {
		t.Fatalf("failed to unmarshal response: %s", err.Error())
	}

	return NewStub(method, req, res)
}

func NewJSONFileStub[I, O protomsg](t *testing.T, stubPath string) Stub {
	data, err := os.ReadFile(stubPath)
	if err != nil {
		t.Fatalf("failed to read stub file %q: %s", stubPath, data)
	}

	var s struct {
		Service string `json:"service"`
		Method  string `json:"method"`
		Input   struct {
			Contains *I `json:"contains"`
		} `json:"input"`
		Output struct {
			Data *O `json:"data"`
		} `json:"output"`
	}

	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("failed to unmarshal stub(%s): %s", stubPath, err.Error())
	}

	return NewStub(s.Method, s.Request, s.Response)
}
*/
