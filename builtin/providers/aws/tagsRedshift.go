package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/hashicorp/terraform/helper/schema"
)

func setTagsRedshift(conn *redshift.Redshift, d *schema.ResourceData, arn string) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTagsRedshift(tagsFromMapRedshift(o), tagsFromMapRedshift(n))

		// Set tags
		if len(remove) > 0 {
			log.Printf("[DEBUG] Removing tags: %#v", remove)
			k := make([]*string, len(remove), len(remove))
			for i, t := range remove {
				k[i] = t.Key
			}

			_, err := conn.DeleteTags(&redshift.DeleteTagsInput{
				ResourceName: aws.String(arn),
				TagKeys:      k,
			})
			if err != nil {
				return err
			}
		}
		if len(create) > 0 {
			log.Printf("[DEBUG] Creating tags: %#v", create)
			_, err := conn.CreateTags(&redshift.CreateTagsInput{
				ResourceName: aws.String(arn),
				Tags:         create,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func diffTagsRedshift(oldTags, newTags []*redshift.Tag) ([]*redshift.Tag, []*redshift.Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[*t.Key] = *t.Value
	}

	// Build the list of what to remove
	var remove []*redshift.Tag
	for _, t := range oldTags {
		old, ok := create[*t.Key]
		if !ok || old != *t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return tagsFromMapRedshift(create), remove
}

func tagsFromMapRedshift(m map[string]interface{}) []*redshift.Tag {
	result := make([]*redshift.Tag, 0, len(m))
	for k, v := range m {
		result = append(result, &redshift.Tag{
			Key:   aws.String(k),
			Value: aws.String(v.(string)),
		})
	}

	return result
}

func tagsToMapRedshift(ts []*redshift.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range ts {
		result[*t.Key] = *t.Value
	}

	return result
}
