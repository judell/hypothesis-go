# Hypothesis client for Go

The [Hypothesis web annotation system](https://web.hypothes.is) comprises a [server](https://github.com/hypothesis/h), a [client](https://github.com/hypothesis/client), and an [API](https://h.readthedocs.io/en/latest/api/). This package provides an interface to the API's [profile](https://h.readthedocs.io/en/latest/api-reference/#tag/profile) and [search](https://h.readthedocs.io/en/latest/api-reference/#tag/annotations/paths/~1search/get) endpoints.

The [Steampipe plugin for Hypothesis](https://github.com/turbot/steampipe-plugin-hypothesis) uses this package to enable [search of public and private annotations](https://hub.steampipe.io/plugins/turbot/hypothesis/tables/hypothesis_search).
