package stemcells_test

const noNewStemcellsAvailableResponse = `
{
  "products": [
    {
      "guid": "p-bosh-4e531084598242b05f9f",
      "deployed_stemcell_version": "3468.13",
      "available_stemcell_versions": [
        "3468.13"
      ]
    }
  ],
  "stemcell_library": [
    {}
  ]
}`

const newStemcellForOneProductResponse = `
{
  "products": [
    {
      "guid": "p-bosh-4e531084598242b05f9f",
      "deployed_stemcell_version": "3468.13",
      "available_stemcell_versions": [
        "3468.16"
      ]
    }
  ]
}`

const newStemcellForOneProductOutput = `
{
  "stemcell_updates": [
	{
	  "stemcell_version": "3468.16",
	  "products": [
		{
		  "product_id": "p-bosh-4e531084598242b05f9f"
		}
	  ]
	}
  ]
}`

const newStemcellsForOneProductResponse = `
{
  "products": [
    {
      "guid": "p-bosh-4e531084598242b05f9f",
      "deployed_stemcell_version": "3468.13",
      "available_stemcell_versions": [
        "3468.15",
        "3468.13",
        "3468.16",
        "3468.14"
      ]
    }
  ]
}`

const newStemcellsForOneProductOutput = `
{
  "stemcell_updates": [
	{
	  "stemcell_version": "3468.16",
	  "products": [
		{
		  "product_id": "p-bosh-4e531084598242b05f9f"
		}
	  ]
	}
  ]
}`

const newStemcellForMultipleProductsResponse = `
{
  "products": [
    {
      "guid": "p-bosh-4e531084598242b05f9f",
      "deployed_stemcell_version": "3468.13",
      "available_stemcell_versions": [
        "3468.16"
      ]
    },
    {
      "guid": "deployed-product-97b88e825c634e430a66",
      "deployed_stemcell_version": "3468.12",
      "available_stemcell_versions": [
        "3468.16"
      ]
    }
  ]
}`

const newStemcellForMultipleProductsOutput = `
{
  "stemcell_updates": [
	{
	  "stemcell_version": "3468.16",
	  "products": [
		{
		  "product_id": "p-bosh-4e531084598242b05f9f"
		},
		{
		  "product_id": "deployed-product-97b88e825c634e430a66"
		}
	  ]
	}
  ]
}`

const newStemcellsForMultipleProductsResponse = `
{
  "products": [
    {
      "guid": "deployed-product-97b88e825c634e430a66",
      "deployed_stemcell_version": "3468.12",
      "available_stemcell_versions": [
        "3468.13",
        "3468.14",
        "3468.16",
        "3468.15"
      ]
    },
    {
      "guid": "p-bosh-4e531084598242b05f9f",
      "deployed_stemcell_version": "3468.13",
      "available_stemcell_versions": [
        "3468.13",
        "3468.17",
        "3468.16"
      ]
    },
    {
      "guid": "la-marsa-beach-28763ba679a87bd",
      "deployed_stemcell_version": "3468.11",
      "available_stemcell_versions": [
        "3468.13",
        "3468.14",
        "3468.16",
        "3468.15"
      ]
    }
  ]
}`

const newStemcellsForMultipleProductsOutput = `{
  "stemcell_updates": [
    {
      "stemcell_version": "3468.16",
      "products": [
        {
          "product_id": "deployed-product-97b88e825c634e430a66"
        },
        {
          "product_id": "la-marsa-beach-28763ba679a87bd"
        }
      ]
    },
    {
      "stemcell_version": "3468.17",
      "products": [
        {
          "product_id": "p-bosh-4e531084598242b05f9f"
        }
      ]
    }
  ]
}`
