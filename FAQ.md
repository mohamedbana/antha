---
layout: default
type: references
navgroup: references
shortname: References
title: FAQ
---

{% include alpha.html %}

*Don't see an answer to your question on here? Ask on the [mailing list](/discuss.html)!*

{% include toc.html %}

## {{site.project_title}}

### Why should I care about this project? {#why}

{{site.project_title}} is the first high level language designed to enable robust, reproducible and composible work in the biosciences. {{site.project_title}} is built on top of established technologies including Go, CouchDB, MORE TEXT HERE.

### Why do we need a high level language for doing biology? {#why}

A high level language enables ease of design, better reproducibility and scalability. Lack of reproducibility is still a major barrier to progress in the biosciences (an article from Amgen reports that they could reproduce findings in only 11% of 53 published papers). 

### What about Biocoder, SBOL, PrPr and Chris? {#why}

None of these are high-level languages capable of incorporating genetic design, experimental design, physical experimental execution and data processing.

### Is {{site.project_title}} production ready? {#readiness}

{{site.project_title}} is currently in "community preview." Many of the pieces that are being integrated into {{site.project_title}} have been in production use at Synthace and other organisations for years. However, the full integration into {{site.project_title_}} is new and will inevitably uncover bugs in the component pieces it has brought together.

### What sort of testing do you do?

{{site.project_title}} uses Chromium's continuous build infrastructure to test
the entire system and each polyfill, individually. See our [build status page](/build/).

### Why choose google GO as the basis?

Go is an open source language designed and built by Google to make building fast, simple and scalable software. Go is capable of concurrency and communicating directly with devices. Learn go [here] (http://go-book.appspot.com/index.html)

## Data Access

### What format does Antha produce data in? {#openData}

{{site.project_title}} translates the proprietary formats of the various devices accessed as part of the process of 
experimeent execution into open JSON based data. The original proprietary format inputs are also captured for archival purposes if
you like those sorts of things.

## Antha Elements

### What sort of things can an element be?

### How small should an element be?

### How do I package a bunch of custom elements together into a workflow? {#packaging}

Antha Elements follow the same packaging structure as GO packages. [Golang](https://code.google.com/p/go-wiki/wiki/PackagePublishing)

### Where do all the low-level details go? Does it matter? 

### Can I integrate scripts from other languages such as Python or Matlab?

Short answer: 

Not yet. But eventually, yes, via the network protocol.

Long answer:
Antha elements are designed to run as microservices communicating via a network using a flow-based approach.
This part of the language is still in very early stage development. Once it is available there will be the option
to integrate existing scripts with the system by wrapping them in middleware designed to convert to and from the
required JSON format and do the other required bits like registering their presence and capabilities to the network.

### How can I contribute to the language?
---
