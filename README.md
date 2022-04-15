# Terraxen Naming Service

Terraxen's naming convention service allows user to create naming conventions and namespaces to keep independant teams stick to resource names for an organization.

# Schema

A schema defines resources and the rules/patterns to use to generate their name. Each schema has version control so that any breaking change automatically creates a new version. This means your digital estate wont change unexpectedly when your schema is updated.

# Namespace

A namespace combines variables with a schema version to give actual values that can be used. Extra parameters can be passed in to the namespace resolution endpoints to create the name.
