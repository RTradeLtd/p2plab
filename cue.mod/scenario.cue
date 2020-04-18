package p2plab

Object :: {
    type: string
    source: string
}

Scenario :: {
    objects: [...Object]
    ...
}

Seed :: {
    neighbors: string
}