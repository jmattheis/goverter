# Use-Cases

Goverter generates boilerplate code for you to convert types to other types.
This can be useful if you have types representing similar data. Use-Cases can
be:

[[toc]]

## Converting between database and API types

Sometimes you have different database and API types because they are:

* Generated with e.g. via [sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc),
  [volatiletech/sqlboiler](https://github.com/volatiletech/sqlboiler)
  [99designs/gqlgen](https://github.com/99designs/gqlgen),
  [protobuf](https://protobuf.dev/getting-started/gotutorial/) or
  [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)
* or because the database ORM requires different typing than the api types
* or because the database models are just different from the API because the
  API must stay backwards compatible

Given any of these reasons it may be useful to generate the conversion via
goverter because it validates that all fields are mapped and reduces the need
for manually converting these types.

## Converting between API versions

When you have to support older API version you could structure you app so it
always uses the latest api types and then convert the latest api types to older
api types via Goverter. The types may have a lot of similarities and therefore,
goverter should be able to automatically convert most of the types.

## Deep copy instances

Goverter allows you to deep copy types by defining a conversion where the
source and target type is the same. This is useful when you want to adjust a
copy of a instance without affecting the source instance.

This can be done at runtime with a library like
[github.com/jinzhu/copier](ttps://github.com/jinzhu/copier) but using generated
code will be much faster and you'll have a compile-time guarantee that the
types can be coverted.
