package main

import (
	"encoding/json"
	"os"

	apex "github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	elt "github.com/aws/aws-sdk-go/service/elastictranscoder"
)

type s3Message struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json: "Name"`
			} `json: "bucket"`
			Object struct {
				Key string `json:key`
			} `json:object`
		} `json: "s3"`
	} `json: "records"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var s3message s3Message
		if err := json.Unmarshal(event, &s3message); err != nil {
			return nil, err
		}
		fileName := s3message.Records[0].S3.Object.Key
		// return os.Getenv("accessKey"), nil
		transcoder := initElasicTranscoder()
		resp := createJob(fileName, transcoder)
		return resp, nil
	})
}

func initElasicTranscoder() *elt.ElasticTranscoder {
	creds := credentials.NewStaticCredentials(os.Getenv("accessKey"), os.Getenv("secretKey"), "")
	awsConfig := aws.NewConfig().WithCredentials(creds).WithRegion("us-west-2")
	svc := elt.New(session.New(awsConfig))
	return svc
}

func createJob(fileName string, transcoder *elt.ElasticTranscoder) string {
	pipeLineID := aws.String("1484274913059-u1b8pl")
	presetID := aws.String("1351620000001-100070")
	inputParams := &elt.JobInput{
		DetectedProperties: nil,
		Encryption:         nil,
		TimeSpan:           nil,
		FrameRate:          aws.String("auto"),
		AspectRatio:        aws.String("auto"),
		Container:          aws.String("mp4"),
		Resolution:         aws.String("auto"),
		Key:                aws.String(fileName),
	}
	ouputParams := &elt.CreateJobOutput{
		Key:      aws.String(fileName + ".mp4"),
		PresetId: presetID,
	}

	jobInput := &elt.CreateJobInput{
		Input:      inputParams,
		Output:     ouputParams,
		PipelineId: pipeLineID,
	}

	_, err := transcoder.CreateJob(jobInput)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return err.Error()
	}
	return "everything worked"
	// fmt.Println(resp)

}
