package p2plab

// Object :: {
//    type: string
//    source: string
//}

object: [Name=_]: {
    type: string
    source: string
}

Scenario :: {
    objects: [...object]
    seed: { ... }
    // enable any fields for benchmark
    benchmark:  { ... }
}