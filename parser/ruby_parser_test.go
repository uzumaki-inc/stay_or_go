package parser_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/uzumaki-inc/StayOrGo/parser"
)

func TestRubyParser_Parse(t *testing.T) {
	t.Parallel()
	// Create a temporary file for testing
	content := `gem 'rails', '~> 6.0'
gem 'nokogiri', git: 'https://self_hosting_git.com/sparklemotion/nokogiri.git'
gem 'puma'`

	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Write content to the file
	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	// Parse the file using RubyParser
	p := parser.RubyParser{}

	libInfoList, err := p.Parse(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	// Assertions
	assert.Len(t, libInfoList, 3)
	assert.Equal(t, "rails", libInfoList[0].Name)
	assert.False(t, libInfoList[0].Skip)
	assert.Equal(t, "nokogiri", libInfoList[1].Name)
	assert.True(t, libInfoList[1].Skip)
	assert.Equal(t, "Not hosted on Github", libInfoList[1].SkipReason)
	assert.Equal(t, "puma", libInfoList[2].Name)
	assert.False(t, libInfoList[2].Skip)
}

func TestRubyParser_GetRepositoryURL(t *testing.T) {
	t.Parallel()
	// Mock HTTP requests for rubygems API
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Setup mock responses
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/rails.json",
		httpmock.NewStringResponder(200, `{"source_code_uri": "https://github.com/rails/rails"}`))

	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/nokogiri.json",
		httpmock.NewStringResponder(200, `{"homepage_uri": "https://github.com/sparklemotion/nokogiri"}`))

	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/puma.json",
		httpmock.NewStringResponder(200, `{"source_code_uri": ""}`))

	// Create initial LibInfo list
	libInfoList := []parser.LibInfo{
		{Name: "rails"},
		{Name: "nokogiri"},
		{Name: "puma"},
	}

	// Run GetRepositoryURL method
	p := parser.RubyParser{}
	updatedLibInfoList := p.GetRepositoryURL(libInfoList)

	// Assertions
	assert.Equal(t, "https://github.com/rails/rails", updatedLibInfoList[0].RepositoryURL)
	assert.Equal(t, "https://github.com/sparklemotion/nokogiri", updatedLibInfoList[1].RepositoryURL)
	assert.Equal(t, "", updatedLibInfoList[2].RepositoryURL)
	assert.True(t, updatedLibInfoList[2].Skip)
	assert.Equal(t, "Does not support libraries hosted outside of Github", updatedLibInfoList[2].SkipReason)
}
