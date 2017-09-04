# certbucket
Certbucket implements `acme/autocert` Cache using a bucket on Google storage. 
If the bucket does not exist, it will be created.

## Usage
```
	cache, err := certbucket.New("your-project-id", "cache-bucket-name")
	if err != nil {
		return err
	} 

	m := autocert.Manager{
		Cache:      cache,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.com"),
	}

	s := &http.Server{
		Addr:      ":https",
		Handler:   router,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	log.Fatal(s.ListenAndServeTLS("", ""))
```

## Testing
Before running tests make sure that your Google Cloud SDK is installed and configured.
Set `GOOGLE_PROJECT_ID` environment variable.

