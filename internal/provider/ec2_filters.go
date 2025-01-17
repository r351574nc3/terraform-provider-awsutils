package provider

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cloudposse/terraform-provider-awsutils/internal/keyvaluetags"
	tfec2 "github.com/cloudposse/terraform-provider-awsutils/internal/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// buildEC2AttributeFilterList takes a flat map of scalar attributes (most
// likely values extracted from a *schema.ResourceData on an EC2-querying
// data source) and produces a []*ec2.Filter representing an exact match
// for each of the given non-empty attributes.
//
// The keys of the given attributes map are the attribute names expected
// by the EC2 API, which are usually either in camelcase or with dash-separated
// words. We conventionally map these to underscore-separated identifiers
// with the same words when presenting these as data source query attributes
// in Terraform.
//
// It's the callers responsibility to transform any non-string values into
// the appropriate string serialization required by the AWS API when
// encoding the given filter. Any attributes given with empty string values
// are ignored, assuming that the user wishes to leave that attribute
// unconstrained while filtering.
//
// The purpose of this function is to create values to pass in
// for the "Filters" attribute on most of the "Describe..." API functions in
// the EC2 API, to aid in the implementation of Terraform data sources that
// retrieve data about EC2 objects.
func buildEC2AttributeFilterList(attrs map[string]string) []*ec2.Filter {
	return tfec2.BuildAttributeFilterList(attrs)
}

// buildEC2TagFilterList takes a []*ec2.Tag and produces a []*ec2.Filter that
// represents exact matches for all of the tag key/value pairs given in
// the tag set.
//
// The purpose of this function is to create values to pass in for
// the "Filters" attribute on most of the "Describe..." API functions
// in the EC2 API, to implement filtering by tag values e.g. in Terraform
// data sources that retrieve data about EC2 objects.
//
// It is conventional for an EC2 data source to include an attribute called
// "tags" which conforms to the schema returned by the tagsSchema() function.
// The value of this can then be converted to a tags slice using tagsFromMap,
// and the result finally passed in to this function.
//
// In Terraform configuration this would then look like this, to constrain
// results by name:
//
// tags {
//   Name = "my-awesome-subnet"
// }
func buildEC2TagFilterList(tags []*ec2.Tag) []*ec2.Filter {
	filters := make([]*ec2.Filter, len(tags))

	for i, tag := range tags {
		filters[i] = &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", *tag.Key)),
			Values: []*string{tag.Value},
		}
	}

	return filters
}

// ec2AttributeFiltersFromMultimap returns an array of EC2 Filter objects to be used when listing resources.
//
// The keys of the specified map are the resource attributes names used in the filter - see the documentation
// for the relevant "Describe" action for a list of the valid names. The resource must match all the filters
// to be included in the result.
// The values of the specified map are lists of resource attribute values used in the filter. The resource can
// match any of the filter values to be included in the result.
// See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html#Filtering_Resources_CLI for more details.
func ec2AttributeFiltersFromMultimap(m map[string][]string) []*ec2.Filter {
	if len(m) == 0 {
		return nil
	}

	filters := []*ec2.Filter{}
	for k, v := range m {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(k),
			Values: aws.StringSlice(v),
		})
	}

	return filters
}

// ec2TagFiltersFromMap returns an array of EC2 Filter objects to be used when listing resources.
//
// The filters represent exact matches for all the resource tags in the given key/value map.
func ec2TagFiltersFromMap(m map[string]interface{}) []*ec2.Filter {
	if len(m) == 0 {
		return nil
	}

	filters := []*ec2.Filter{}
	for _, tag := range keyvaluetags.New(m).IgnoreAws().Ec2Tags() {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", aws.StringValue(tag.Key))),
			Values: []*string{tag.Value},
		})
	}

	return filters
}

// ec2CustomFiltersSchema returns a *schema.Schema that represents
// a set of custom filtering criteria that a user can specify as input
// to a data source that wraps one of the many "Describe..." API calls
// in the EC2 API.
//
// It is conventional for an attribute of this type to be included
// as a top-level attribute called "filter". This is the "catch all" for
// filter combinations that are not possible to express using scalar
// attributes or tags. In Terraform configuration, the custom filter blocks
// then look like this:
//
// filter {
//   name   = "availabilityZone"
//   values = ["us-west-2a", "us-west-2b"]
// }
func ec2CustomFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"values": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// buildEC2CustomFilterList takes the set value extracted from a schema
// attribute conforming to the schema returned by ec2CustomFiltersSchema,
// and transforms it into a []*ec2.Filter representing the same filter
// expressions which is ready to pass into the "Filters" attribute on most
// of the "Describe..." functions in the EC2 API.
//
// This function is intended only to be used in conjunction with
// ec2CustomFitlersSchema. See the docs on that function for more details
// on the configuration pattern this is intended to support.
func buildEC2CustomFilterList(filterSet *schema.Set) []*ec2.Filter {
	if filterSet == nil {
		return []*ec2.Filter{}
	}

	customFilters := filterSet.List()
	filters := make([]*ec2.Filter, len(customFilters))

	for filterIdx, customFilterI := range customFilters {
		customFilterMapI := customFilterI.(map[string]interface{})
		name := customFilterMapI["name"].(string)
		valuesI := customFilterMapI["values"].(*schema.Set).List()
		values := make([]*string, len(valuesI))
		for valueIdx, valueI := range valuesI {
			values[valueIdx] = aws.String(valueI.(string))
		}

		filters[filterIdx] = &ec2.Filter{
			Name:   &name,
			Values: values,
		}
	}

	return filters
}
