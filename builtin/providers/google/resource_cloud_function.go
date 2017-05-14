package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudfunctions/v1beta2"
)

func resourceCloudFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFunctionCreate,
		Read:   resourceCloudFunctionRead,
		Update: resourceCloudFunctionUpdate,
		Delete: resourceCloudFunctionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	function := &cloudfunctions.CloudFunction{
		AvailableMemoryMb: d.Get("available_memory_mb").(int),
		EntryPoint:        d.Get("entry_point").(string),
		EventTrigger: &cloudfunctions.EventTrigger{
			EventType:       "",
			ForceSendFields: []string{},
			NullFields:      []string{},
			Resource:        "",
		},
		HttpsTrigger: &cloudfunctions.HTTPSTrigger{
			ForceSendFields: []string{},
			NullFields:      []string{},
			Url:             "",
		},
		Name:             d.Get("name").(string), // match pattern `projects/*/locations/*/functions/*`
		SourceArchiveUrl: d.Get("source_archive_url").(string),
		SourceRepository: &cloudfunctions.SourceRepository{
			Branch:           "",
			DeployedRevision: "",
			ForceSendFields:  []string{},
			NullFields:       []string{},
			RepositoryUrl:    "",
			Revision:         "",
			SourcePath:       "",
			Tag:              "",
		},
		Timeout: d.Get("timeout").(string),
	}

	call := config.clientFunctions.Projects.Locations.Functions.Create(location, function)
	res, err := call.Do()
	if err != nil {
		return err
	}

	d.SetId(res.Name)

	return nil
}

func resourceCloudFunctionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientFunctions.Projects.Locations.Functions.Get(name)
	out, err := call.Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Pubsub Topic %q", name))
	}

	return nil
}

func resourceCloudFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientPubsub.Projects.Topics.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err
	}

	return nil
}
