package manifestlist

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/manifest/ocischema"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const expectedManifestListSerialization = `{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
   "manifests": [
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "digest": "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
         "size": 985,
         "platform": {
            "architecture": "amd64",
            "os": "linux",
            "features": [
               "sse4"
            ]
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "digest": "sha256:6346340964309634683409684360934680934608934608934608934068934608",
         "size": 2392,
         "platform": {
            "architecture": "sun4m",
            "os": "sunos"
         }
      }
   ]
}`

func makeTestManifestList(t *testing.T, mediaType string) ([]ManifestDescriptor, *DeserializedManifestList) {
	manifestDescriptors := []ManifestDescriptor{
		{
			Descriptor: distribution.Descriptor{
				MediaType: "application/vnd.docker.distribution.manifest.v2+json",
				Digest:    "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
				Size:      985,
			},
			Platform: PlatformSpec{
				Architecture: "amd64",
				OS:           "linux",
				Features:     []string{"sse4"},
			},
		},
		{
			Descriptor: distribution.Descriptor{
				MediaType: "application/vnd.docker.distribution.manifest.v2+json",
				Digest:    "sha256:6346340964309634683409684360934680934608934608934608934068934608",
				Size:      2392,
			},
			Platform: PlatformSpec{
				Architecture: "sun4m",
				OS:           "sunos",
			},
		},
	}

	deserialized, err := FromDescriptorsWithMediaType(manifestDescriptors, mediaType)
	if err != nil {
		t.Fatalf("error creating DeserializedManifestList: %v", err)
	}

	return manifestDescriptors, deserialized
}

func TestManifestList(t *testing.T) {
	manifestDescriptors, deserialized := makeTestManifestList(t, MediaTypeManifestList)
	mediaType, canonical, _ := deserialized.Payload()

	if mediaType != MediaTypeManifestList {
		t.Fatalf("unexpected media type: %s", mediaType)
	}

	// Check that the canonical field is the same as json.MarshalIndent
	// with these parameters.
	expected, err := json.MarshalIndent(&deserialized.ManifestList, "", "   ")
	if err != nil {
		t.Fatalf("error marshaling manifest list: %v", err)
	}
	if !bytes.Equal(expected, canonical) {
		t.Fatalf("manifest bytes not equal:\nexpected:\n%s\nactual:\n%s\n", string(expected), string(canonical))
	}

	// Check that the canonical field has the expected value.
	if !bytes.Equal([]byte(expectedManifestListSerialization), canonical) {
		t.Fatalf("manifest bytes not equal:\nexpected:\n%s\nactual:\n%s\n", expectedManifestListSerialization, string(canonical))
	}

	var unmarshalled DeserializedManifestList
	if err := json.Unmarshal(deserialized.canonical, &unmarshalled); err != nil {
		t.Fatalf("error unmarshaling manifest: %v", err)
	}

	if !reflect.DeepEqual(&unmarshalled, deserialized) {
		t.Fatalf("manifests are different after unmarshaling: %v != %v", unmarshalled, *deserialized)
	}

	references := deserialized.References()
	if len(references) != 2 {
		t.Fatalf("unexpected number of references: %d", len(references))
	}
	for i := range references {
		platform := manifestDescriptors[i].Platform
		expectedPlatform := &v1.Platform{
			Architecture: platform.Architecture,
			OS:           platform.OS,
			OSFeatures:   platform.OSFeatures,
			OSVersion:    platform.OSVersion,
			Variant:      platform.Variant,
		}
		if !reflect.DeepEqual(references[i].Platform, expectedPlatform) {
			t.Fatalf("unexpected value %d returned by References: %v", i, references[i])
		}
		references[i].Platform = nil
		if !reflect.DeepEqual(references[i], manifestDescriptors[i].Descriptor) {
			t.Fatalf("unexpected value %d returned by References: %v", i, references[i])
		}
	}
}

// TODO (mikebrow): add annotations on the manifest list (index) and support for
// empty platform structs (move to Platform *Platform `json:"platform,omitempty"`
// from current Platform PlatformSpec `json:"platform"`) in the manifest descriptor.
// Requires changes to distribution/distribution/manifest/manifestlist.ManifestList and .ManifestDescriptor
// and associated serialization APIs in manifestlist.go. Or split the OCI index and
// docker manifest list implementations, which would require a lot of refactoring.
const expectedOCIImageIndexSerialization = `{
   "schemaVersion": 2,
   "mediaType": "application/vnd.oci.image.index.v1+json",
   "manifests": [
      {
         "mediaType": "application/vnd.oci.image.manifest.v1+json",
         "digest": "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
         "size": 985,
         "platform": {
            "architecture": "amd64",
            "os": "linux",
            "features": [
               "sse4"
            ]
         }
      },
      {
         "mediaType": "application/vnd.oci.image.manifest.v1+json",
         "digest": "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
         "size": 985,
         "annotations": {
            "platform": "none"
         },
         "platform": {
            "architecture": "",
            "os": ""
         }
      },
      {
         "mediaType": "application/vnd.oci.image.manifest.v1+json",
         "digest": "sha256:6346340964309634683409684360934680934608934608934608934068934608",
         "size": 2392,
         "annotations": {
            "what": "for"
         },
         "platform": {
            "architecture": "sun4m",
            "os": "sunos"
         }
      }
   ]
}`

func makeTestOCIImageIndex(t *testing.T, mediaType string) ([]ManifestDescriptor, *DeserializedManifestList) {
	manifestDescriptors := []ManifestDescriptor{
		{
			Descriptor: distribution.Descriptor{
				MediaType: "application/vnd.oci.image.manifest.v1+json",
				Digest:    "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
				Size:      985,
			},
			Platform: PlatformSpec{
				Architecture: "amd64",
				OS:           "linux",
				Features:     []string{"sse4"},
			},
		},
		{
			Descriptor: distribution.Descriptor{
				MediaType:   "application/vnd.oci.image.manifest.v1+json",
				Digest:      "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
				Size:        985,
				Annotations: map[string]string{"platform": "none"},
			},
		},
		{
			Descriptor: distribution.Descriptor{
				MediaType:   "application/vnd.oci.image.manifest.v1+json",
				Digest:      "sha256:6346340964309634683409684360934680934608934608934608934068934608",
				Size:        2392,
				Annotations: map[string]string{"what": "for"},
			},
			Platform: PlatformSpec{
				Architecture: "sun4m",
				OS:           "sunos",
			},
		},
	}

	deserialized, err := FromDescriptorsWithMediaType(manifestDescriptors, mediaType)
	if err != nil {
		t.Fatalf("error creating DeserializedManifestList: %v", err)
	}

	return manifestDescriptors, deserialized
}

func TestOCIImageIndex(t *testing.T) {
	manifestDescriptors, deserialized := makeTestOCIImageIndex(t, v1.MediaTypeImageIndex)

	mediaType, canonical, _ := deserialized.Payload()

	if mediaType != v1.MediaTypeImageIndex {
		t.Fatalf("unexpected media type: %s", mediaType)
	}

	// Check that the canonical field is the same as json.MarshalIndent
	// with these parameters.
	expected, err := json.MarshalIndent(&deserialized.ManifestList, "", "   ")
	if err != nil {
		t.Fatalf("error marshaling manifest list: %v", err)
	}
	if !bytes.Equal(expected, canonical) {
		t.Fatalf("manifest bytes not equal:\nexpected:\n%s\nactual:\n%s\n", string(expected), string(canonical))
	}

	// Check that the canonical field has the expected value.
	if !bytes.Equal([]byte(expectedOCIImageIndexSerialization), canonical) {
		t.Fatalf("manifest bytes not equal:\nexpected:\n%s\nactual:\n%s\n", expectedOCIImageIndexSerialization, string(canonical))
	}

	var unmarshalled DeserializedManifestList
	if err := json.Unmarshal(deserialized.canonical, &unmarshalled); err != nil {
		t.Fatalf("error unmarshaling manifest: %v", err)
	}

	if !reflect.DeepEqual(&unmarshalled, deserialized) {
		t.Fatalf("manifests are different after unmarshaling: %v != %v", unmarshalled, *deserialized)
	}

	references := deserialized.References()
	if len(references) != 3 {
		t.Fatalf("unexpected number of references: %d", len(references))
	}
	for i := range references {
		platform := manifestDescriptors[i].Platform
		expectedPlatform := &v1.Platform{
			Architecture: platform.Architecture,
			OS:           platform.OS,
			OSFeatures:   platform.OSFeatures,
			OSVersion:    platform.OSVersion,
			Variant:      platform.Variant,
		}
		if !reflect.DeepEqual(references[i].Platform, expectedPlatform) {
			t.Fatalf("unexpected value %d returned by References: %v", i, references[i])
		}
		references[i].Platform = nil
		if !reflect.DeepEqual(references[i], manifestDescriptors[i].Descriptor) {
			t.Fatalf("unexpected value %d returned by References: %v", i, references[i])
		}
	}
}

func mediaTypeTest(t *testing.T, contentType string, mediaType string, shouldError bool) {
	var m *DeserializedManifestList
	if contentType == MediaTypeManifestList {
		_, m = makeTestManifestList(t, mediaType)
	} else {
		_, m = makeTestOCIImageIndex(t, mediaType)
	}

	_, canonical, err := m.Payload()
	if err != nil {
		t.Fatalf("error getting payload, %v", err)
	}

	unmarshalled, descriptor, err := distribution.UnmarshalManifest(
		contentType,
		canonical)

	if shouldError {
		if err == nil {
			t.Fatalf("bad content type should have produced error")
		}
	} else {
		if err != nil {
			t.Fatalf("error unmarshaling manifest, %v", err)
		}

		asManifest := unmarshalled.(*DeserializedManifestList)
		if asManifest.MediaType != mediaType {
			t.Fatalf("Bad media type '%v' as unmarshalled", asManifest.MediaType)
		}

		if descriptor.MediaType != contentType {
			t.Fatalf("Bad media type '%v' for descriptor", descriptor.MediaType)
		}

		unmarshalledMediaType, _, _ := unmarshalled.Payload()
		if unmarshalledMediaType != contentType {
			t.Fatalf("Bad media type '%v' for payload", unmarshalledMediaType)
		}
	}
}

func TestMediaTypes(t *testing.T) {
	mediaTypeTest(t, MediaTypeManifestList, "", true)
	mediaTypeTest(t, MediaTypeManifestList, MediaTypeManifestList, false)
	mediaTypeTest(t, MediaTypeManifestList, MediaTypeManifestList+"XXX", true)
	mediaTypeTest(t, v1.MediaTypeImageIndex, "", false)
	mediaTypeTest(t, v1.MediaTypeImageIndex, v1.MediaTypeImageIndex, false)
	mediaTypeTest(t, v1.MediaTypeImageIndex, v1.MediaTypeImageIndex+"XXX", true)
}

func TestValidateManifest(t *testing.T) {
	manifest := ocischema.Manifest{
		Config: distribution.Descriptor{Size: 1},
		Layers: []distribution.Descriptor{{Size: 2}},
	}
	index := ManifestList{
		Manifests: []ManifestDescriptor{
			{Descriptor: distribution.Descriptor{Size: 3}},
		},
	}
	t.Run("valid", func(t *testing.T) {
		b, err := json.Marshal(index)
		if err != nil {
			t.Fatal("unexpected error marshaling index", err)
		}
		if err := validateIndex(b); err != nil {
			t.Error("index should be valid", err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		b, err := json.Marshal(manifest)
		if err != nil {
			t.Fatal("unexpected error marshaling manifest", err)
		}
		if err := validateIndex(b); err == nil {
			t.Error("manifest should not be valid")
		}
	})
}
