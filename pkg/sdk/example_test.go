package sdk

import "fmt"

func ExampleNewClient() {
	client := NewClient("https://api.example.com", "api-key")
	fmt.Println(client != nil)
	// Output: true
}

func ExampleVolumeAttachmentInput() {
	attachment := VolumeAttachmentInput{
		VolumeID:  "vol-123",
		MountPath: "/data",
	}
	fmt.Printf("%s:%s\n", attachment.VolumeID, attachment.MountPath)
	// Output: vol-123:/data
}

func ExampleLifecycleRule() {
	rule := LifecycleRule{
		BucketName:     "logs",
		Prefix:         "app/",
		ExpirationDays: 30,
		Enabled:        true,
	}
	fmt.Printf("%s:%s:%d:%t\n", rule.BucketName, rule.Prefix, rule.ExpirationDays, rule.Enabled)
	// Output: logs:app/:30:true
}
