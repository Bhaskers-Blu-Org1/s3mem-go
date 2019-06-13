/*
###############################################################################
# Licensed Materials - Property of IBM Copyright IBM Corporation 2017, 2019. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP
# Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################
*/

package s3mem

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.ibm.com/open-razee/s3mem-go/s3mem/s3memerr"
)

const opGetBucketVersioning = "GetBucketVersioning"

//GetBucketVersioningRequest ...
func (c *Client) GetBucketVersioningRequest(input *s3.GetBucketVersioningInput) s3.GetBucketVersioningRequest {
	if input == nil {
		input = &s3.GetBucketVersioningInput{}
	}
	output := &s3.GetBucketVersioningOutput{}
	op := &aws.Operation{
		Name:       opGetBucketVersioning,
		HTTPMethod: "GET",
		HTTPPath:   "/{Bucket}?versioning",
	}
	req := c.NewRequest(op, input, output)
	return s3.GetBucketVersioningRequest{Request: req, Input: input}
}

func getBucketVersioning(req *aws.Request) {
	if !IsBucketExist(req.Params.(*s3.GetBucketVersioningInput).Bucket) {
		req.Error = s3memerr.NewError(s3.ErrCodeNoSuchBucket, "", nil, req.Params.(*s3.GetBucketVersioningInput).Bucket, nil, nil)
		return
	}
	_, obj := GetBucketVersioning(req.Params.(*s3.GetBucketVersioningInput).Bucket)
	if obj == nil {
		req.Error = s3memerr.NewError(s3.ErrCodeNoSuchBucket, "", nil, req.Params.(*s3.GetBucketVersioningInput).Bucket, nil, nil)
		return
	}
	switch obj.MFADelete {
	case s3.MFADeleteEnabled:
		req.Data.(*s3.GetBucketVersioningOutput).MFADelete = s3.MFADeleteStatusEnabled
	case s3.MFADeleteDisabled:
		req.Data.(*s3.GetBucketVersioningOutput).MFADelete = s3.MFADeleteStatusDisabled
	default:
	}
	req.Data.(*s3.GetBucketVersioningOutput).Status = obj.Status
	//This is needed just to logResponse when requested
	body, _ := json.MarshalIndent(req.Data.(*s3.GetBucketVersioningOutput), "", "  ")
	req.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader(body))
}