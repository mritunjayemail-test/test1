package common

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/packer/template/interpolate"
)

// AccessConfig is for common configuration related to AWS access
type AccessConfig struct {
	/* The access key used to communicate with AWS. [Learn
	how to set this.](/docs/builders/amazon.html#specifying-amazon-credentials)
	*/
	AccessKey string `mapstructure:"access_key"`
	/* This option is useful if you use
	another cloud provider that provide a compatible API with aws EC2.
	specify another endpoint like this "<https://ec2.another.endpoint>..com"
	*/
	CustomEndpointEc2 string `mapstructure:"custom_endpoint_ec2"`
	/* The MFA [TOTP](https://en.wikipedia.org/wiki/Time-based_One-time_Password_Algorithm)
	code. This should probably be a user variable since it changes all the time.
	*/
	MFACode string `mapstructure:"mfa_code"`
	/* The profile to use in the shared credentials file for
	AWS. See Amazon's documentation on [specifying
	profiles](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-profiles)
	for more details.
	*/
	ProfileName string `mapstructure:"profile"`
	/* The name of the region, such as `us-east-1`, in which to
	launch the EC2 instance to create the AMI.
	*/
	RawRegion string `mapstructure:"region"`
	/* Secret key used to communicate with AWS. [Learn
	how to set this.](/docs/builders/amazon.html#specifying-amazon-credentials)
	*/
	SecretKey string `mapstructure:"secret_key"`
	/* Set to true if you want to skip
	validation of the region configuration option. Default `false`.
	*/
	SkipValidation bool `mapstructure:"skip_region_validation"`
	/* The access token to use. This is different from the
	access key and secret key. If you're not sure what this is, then you
	probably don't need it. This will also be read from the `AWS_SESSION_TOKEN`
	environmental variable.
	*/
	Token   string `mapstructure:"token"`
	session *session.Session
}

// Config returns a valid aws.Config object for access to AWS services, or
// an error if the authentication and region couldn't be resolved
func (c *AccessConfig) Session() (*session.Session, error) {
	if c.session != nil {
		return c.session, nil
	}

	region, err := c.Region()
	if err != nil {
		return nil, err
	}

	if c.ProfileName != "" {
		if err := os.Setenv("AWS_PROFILE", c.ProfileName); err != nil {
			log.Printf("Set env error: %s", err)
		}
	}

	config := aws.NewConfig().WithRegion(region).WithMaxRetries(11).WithCredentialsChainVerboseErrors(true)

	if c.CustomEndpointEc2 != "" {
		config = config.WithEndpoint(c.CustomEndpointEc2)
	}

	if c.AccessKey != "" {
		creds := credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     c.AccessKey,
						SecretAccessKey: c.SecretKey,
						SessionToken:    c.Token,
					},
				},
			})
		config = config.WithCredentials(creds)
	}

	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	}
	if c.MFACode != "" {
		opts.AssumeRoleTokenProvider = func() (string, error) {
			return c.MFACode, nil
		}
	}
	c.session, err = session.NewSessionWithOptions(opts)
	if err != nil {
		return nil, err
	}

	return c.session, nil
}

// Region returns the aws.Region object for access to AWS services, requesting
// the region from the instance metadata if possible.
func (c *AccessConfig) Region() (string, error) {
	if c.RawRegion != "" {
		if !c.SkipValidation {
			if valid := ValidateRegion(c.RawRegion); !valid {
				return "", fmt.Errorf("Not a valid region: %s", c.RawRegion)
			}
		}
		return c.RawRegion, nil
	}

	sess := session.New()
	ec2meta := ec2metadata.New(sess)
	identity, err := ec2meta.GetInstanceIdentityDocument()
	if err != nil {
		log.Println("Error getting region from metadata service, "+
			"probably because we're not running on AWS.", err)
		return "", nil
	}
	return identity.Region, nil
}

func (c *AccessConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if c.RawRegion != "" && !c.SkipValidation {
		if valid := ValidateRegion(c.RawRegion); !valid {
			errs = append(errs, fmt.Errorf("Unknown region: %s", c.RawRegion))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
