# Filters

Whenever an Event occurs, all of the Filters you defined in a Chain get executed with the data the Event provided. Only if the Event passes all of the Filters, the configured Actions in this Chain get then executed.

At the moment, there are 2 different types of filters:

+ temmplate
+ starlark

## Template

"Template" filter type uses [text/template](https://golang.org/pkg/text/template/) Go package to execute filters. Beehive exposes in the template a few helper functions and most of the [strings](https://golang.org/pkg/strings/) package functions to make templates more powerful. Also an important thing is that the template for filters should go inside `{{test ...}}`. That's pretty much it.

For example, let's check if `text` contains word "beehive":

```clojure
{{test Contains (ToLower .text) "beehive"}}
```

See [Beehive wiki](https://github.com/muesli/beehive/wiki/Filters) for more information.

## Starlark

[Starlark](https://github.com/bazelbuild/starlark) is a dialect of Python created for [Bazel](https://bazel.build/) configuration files. If you know Python, you already know Starlark. To make sure if a syntax feature is supported check [the specification](https://github.com/google/starlark-go/blob/master/doc/spec.md).

Beehive executes `main` function from the filter, passing all variables inside as keyword arguments. The function must return a boolean result.

For example, let's check if title or description of an RSS feed contains any of the given terms:

```python
def main(title, description, **kwargs):
    text = title + " " + description
    text = text.lower()
    terms = ["beehive", "muesli"]
    for term in terms:
        if term in text:
        return True
    return False
```

"Starlark" filters are more verbose that "template" but much more powerful and turing-complete.
