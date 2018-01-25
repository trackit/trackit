package s3

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
)

// dumbS3Manager implements custom request and signing code in order to
// circumvent a limitation of the AWS SDK described here:
//   https://stackoverflow.com/questions/48285381/how-to-use-getobject-when-key-has-leading-slashes
type dumbS3Manager struct {
	creds  *credentials.Credentials
	signer *v4.Signer
}

// readToWriter reads everything from a reader and writes it to a writer.
func readToWriter(reader io.Reader, writer io.Writer) error {
	var bufReader = bufio.NewReader(reader)
	var bufWriter = bufio.NewWriter(writer)
	_, err := bufWriter.ReadFrom(bufReader)
	return err
}

// init initializes the dumbS3Manager's signer from a set of credentials.
func (dm *dumbS3Manager) init(creds *credentials.Credentials) {
	dm.signer = v4.NewSigner(creds)
}

// rawS3GetObjectToBuffer reads an S3 object to a buffer.
func (dm dumbS3Manager) rawS3GetObjectToBuffer(
	ctx context.Context, client *http.Client, region, bucket, key string,
) ([]byte, error) {
	body, err := dm.rawS3GetObjectToReader(ctx, client, region, bucket, key)
	if err != nil {
		return nil, err
	}
	bufWrapper := bytes.NewBuffer(make([]byte, 0, 0x8000))
	err = readToWriter(body, bufWrapper)
	if err != nil {
		return nil, err
	}
	return bufWrapper.Bytes(), nil
}

// rawS3GetObjectToReader returns the reader for an S3 object.
func (dm dumbS3Manager) rawS3GetObjectToReader(ctx context.Context, client *http.Client, region, bucket, key string) (io.ReadCloser, error) {
	req, err := dm.rawS3GetObject(region, bucket, key)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

// rawS3GetObject returns the signed requests which will request the given S3
// object.
func (dm dumbS3Manager) rawS3GetObject(region, bucket, key string) (*http.Request, error) {
	url := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		bucket,
		region,
		rest.EscapePath(key, false),
	)
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}
	_, err = dm.signer.Sign(request, nil, "s3", region, time.Now())
	if err != nil {
		return nil, err
	}
	return request, nil
}
