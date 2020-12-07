# acme-inc-shop

The acme-inc-shop is an example for organizing multiple test-cases around the same online shop.

Note that this is not the only way how to structure this and this example may change over time.
We welcome feedback from your experience.

## Context

A typical online shop resolves around a few common pages:

1. Landing pages (including the start page, but also certain ad campaign specific ones)
1. Product pages
1. Basket page

These are normally the most important and most frequent pages. If any page breaks, the core business process for customers to submit orders is broken.

## Structure

1. `scenarios/` contains the different user scenarios to describe the visitor behavior.
1. `components/` contains helper functions to act on the different pages or interact with the API
1. `*.mjs` The top level javascript files describe the different test cases

In this example the scenarios should be a high level description on the common visiot behavior, e.g. "visits the landing page, visits 2-5 product pages, adds 1 product to the basket, performs a checkout". To implement this high level description, the components provide the necessary functionality to interact with the shop.

The top level javascript files are left with the general test case setup:

1. Configuration (URL of the shop, authentication, ...)
1. Arrival Phases (How many visitors to start)
1. Test Configuration (Target Sites, Region, Cluster Sizing)
1. Session Weights ()
