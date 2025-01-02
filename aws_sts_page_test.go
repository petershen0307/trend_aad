package trendaad

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StsParsing(t *testing.T) {
	jsonStr := `{
    "000568857918": {
        "AAD-READONLY_000568857918": "An error occurred (AccessDenied) when calling the AssumeRoleWithSAML\noperation:\nNot\nauthorized to perform sts:AssumeRoleWithSAML"
    },
    "147995105371": {
        "AAD-READONLY_147995105371": {
            "AccessKeyId": "ASASE5JSDBN2R7C7S25",
            "SecretAccessKey": "aWD5cM4HNTVgbGtmS2HVfgwNwniY2+epDB/yUkJE",
            "SessionToken": "IQoJb3JpZ2luX2VjEPb//////////7Bg0mG1IvFfvRo5P0y60vC6PC8qAIhAKCJjDqdnPR79gOpy0U7zbDL4KyI",
            "Expiration": "2025-01-02T06:53:16+00:00"
        },
        "name": "c1ws-prod-us"
    },
    "658949957415": {
        "AAD-RD_658949957415": {
            "AccessKeyId": "SIAZS3D2M4T45GPHKXJ",
            "SecretAccessKey": "04UgGeQAJBwk264YdR6SLqCt52W/22C9jWWPmOJc",
            "SessionToken": "IQoJb3JpZ2luX2VjEPbwEaCXVzLWVhc3QtMSJIMEYCIQD",
            "Expiration": "2025-01-02T17:53:18+00:00"
        },
        "name": "xdr-sae-stg",
        "AAD-Contributor_658949957415": {
            "AccessKeyId": "ASIAZS3D2M4TS4ME7WY",
            "SecretAccessKey": "wLkICoCrwhPzgFLVWRRDzd3t047BWZEu4gqHoh+u",
            "SessionToken": "IQoJb3JpZ2luX2VjEPb//////////wEaCXVzLWVhc3QtMSJIMEYCIQD+3",
            "Expiration": "2025-01-02T06:53:23+00:00"
        }
    }
}`
	sts := TrendAwsSts{}
	assert.NoError(t, json.Unmarshal([]byte(jsonStr), &sts))
	assert.Len(t, sts, 3)
	assert.Contains(t, sts, "000568857918")
	assert.Len(t, sts["000568857918"], 1)
	assert.Contains(t, sts["000568857918"], "AAD-READONLY_000568857918")
	assert.Equal(t, "An error occurred (AccessDenied) when calling the AssumeRoleWithSAML\noperation:\nNot\nauthorized to perform sts:AssumeRoleWithSAML", sts["000568857918"]["AAD-READONLY_000568857918"].(string))

	assert.Contains(t, sts, "147995105371")
	assert.Len(t, sts["147995105371"], 2)
	assert.Contains(t, sts["147995105371"], "AAD-READONLY_147995105371")
	assert.Equal(t, sts["147995105371"]["name"].(string), "c1ws-prod-us")
	assert.Equal(t, "ASASE5JSDBN2R7C7S25", sts["147995105371"]["AAD-READONLY_147995105371"].(map[string]interface{})["AccessKeyId"].(string))

	assert.Contains(t, sts, "658949957415")
}

func Test_flushing(t *testing.T) {
	jsonStr := `{
		"000568857918": {
			"AAD-READONLY_000568857918": "An error occurred (AccessDenied) when calling the AssumeRoleWithSAML\noperation:\nNot\nauthorized to perform sts:AssumeRoleWithSAML"
		},
		"147995105371": {
			"AAD-READONLY_147995105371": {
				"AccessKeyId": "ASASE5JSDBN2R7C7S25",
				"SecretAccessKey": "aWD5cM4HNTVgbGtmS2HVfgwNwniY2+epDB/yUkJE",
				"SessionToken": "IQoJb3JpZ2luX2VjEPb//////////7Bg0mG1IvFfvRo5P0y60vC6PC8qAIhAKCJjDqdnPR79gOpy0U7zbDL4KyI",
				"Expiration": "2025-01-02T06:53:16+00:00"
			},
			"name": "c1ws-prod-us"
		},
		"658949957415": {
			"AAD-RD_658949957415": {
				"AccessKeyId": "SIAZS3D2M4T45GPHKXJ",
				"SecretAccessKey": "04UgGeQAJBwk264YdR6SLqCt52W/22C9jWWPmOJc",
				"SessionToken": "IQoJb3JpZ2luX2VjEPbwEaCXVzLWVhc3QtMSJIMEYCIQD",
				"Expiration": "2025-01-02T17:53:18+00:00"
			},
			"name": "xdr-sae-stg",
			"AAD-Contributor_658949957415": {
				"AccessKeyId": "ASIAZS3D2M4TS4ME7WY",
				"SecretAccessKey": "wLkICoCrwhPzgFLVWRRDzd3t047BWZEu4gqHoh+u",
				"SessionToken": "IQoJb3JpZ2luX2VjEPb//////////wEaCXVzLWVhc3QtMSJIMEYCIQD+3",
				"Expiration": "2025-01-02T06:53:23+00:00"
			}
		}
	}`
	sts := TrendAwsSts{}
	assert.NoError(t, json.Unmarshal([]byte(jsonStr), &sts))
	buffer := bytes.NewBuffer([]byte{})
	sts.FlushAwsCredential(buffer)
	output := buffer.String()
	assert.Contains(t, output, "[c1ws-prod-us_AAD-READONLY_147995105371]\n")
	assert.Contains(t, output, "aws_access_key_id = ASASE5JSDBN2R7C7S25\n")
	assert.Contains(t, output, "aws_secret_access_key = aWD5cM4HNTVgbGtmS2HVfgwNwniY2+epDB/yUkJE\n")
	assert.Contains(t, output, "aws_session_token = IQoJb3JpZ2luX2VjEPb//////////7Bg0mG1IvFfvRo5P0y60vC6PC8qAIhAKCJjDqdnPR79gOpy0U7zbDL4KyI\n")
	assert.Contains(t, output, "[xdr-sae-stg_AAD-Contributor_658949957415]\n")
	assert.Contains(t, output, "aws_access_key_id = ASIAZS3D2M4TS4ME7WY\n")
	assert.Contains(t, output, "aws_secret_access_key = wLkICoCrwhPzgFLVWRRDzd3t047BWZEu4gqHoh+u\n")
	assert.Contains(t, output, "aws_session_token = IQoJb3JpZ2luX2VjEPb//////////wEaCXVzLWVhc3QtMSJIMEYCIQD+3\n")
	assert.Contains(t, output, "[xdr-sae-stg_AAD-RD_658949957415]\n")
	assert.Contains(t, output, "aws_access_key_id = SIAZS3D2M4T45GPHKXJ\n")
	assert.Contains(t, output, "aws_secret_access_key = 04UgGeQAJBwk264YdR6SLqCt52W/22C9jWWPmOJc\n")
	assert.Contains(t, output, "aws_session_token = IQoJb3JpZ2luX2VjEPbwEaCXVzLWVhc3QtMSJIMEYCIQD\n")
}

func Test_flushing_specialStructure(t *testing.T) {
	jsonStr := `{
    "761533865683": {
        "AAD-ADMIN_761533865683": {
            "AccessKeyId": "ASIA3CTXUT3JYCE52AFZ",
            "SecretAccessKey": "u0EZ2er+",
            "SessionToken": "IQoJb3JpZ2luX2VjEPb//////////wEaCXVzLWVhc3QtMSJHMEUCIQDsD9VvkdYvOLsXwyTba1tms5ZP177Fso",
            "Expiration": "2025-01-02T06:53:26+00:00"
        },
        "name": "xdr-crp-prod",
        "AAD-OPS_761533865683": "An error occurred (AccessDenied) when calling the AssumeRoleWithSAML operation:\nNot\nauthorized\nto perform sts:AssumeRoleWithSAML"
    }
}`
	sts := TrendAwsSts{}
	assert.NoError(t, json.Unmarshal([]byte(jsonStr), &sts))
	buffer := bytes.NewBuffer([]byte{})
	sts.FlushAwsCredential(buffer)
	output := buffer.String()
	assert.Equal(t, `[xdr-crp-prod_AAD-ADMIN_761533865683]
aws_access_key_id = ASIA3CTXUT3JYCE52AFZ
aws_secret_access_key = u0EZ2er+
aws_session_token = IQoJb3JpZ2luX2VjEPb//////////wEaCXVzLWVhc3QtMSJHMEUCIQDsD9VvkdYvOLsXwyTba1tms5ZP177Fso
`, output)
}
