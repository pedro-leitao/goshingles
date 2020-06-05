# An idiomatic implementation of word shingles (ngrams) in Golang #

Goshingles is an implementation of word n-grams in Go. It has a (relatively) small memory footprint, and is fast. It can ingest just about any text form, and it will try its best to split the content into sentences and words from which to compute n-grams.

Here's an example of use:

    import shingles
    (...)
    var sh Shingles
    var text string = `
        Everything, everything, everything, everything..
        In its right place
        In its right place
        In its right place
        Right place

        Yesterday I woke up sucking a lemon
        Yesterday I woke up sucking a lemon
        Yesterday I woke up sucking a lemon
        Yesterday I woke up sucking a lemon

        Everything, everything, everything..
        In its right place
        In its right place
        Right place

        There are two colours in my head
        There are two colours in my head
        What is that you tried to say?
        What was that you tried to say?
        Tried to say.. tried to say..
        Tried to say.. tried to say..

        Everything in its right place `

    sh.Initialize(TRIGRAM)
    sh.Incorporate(text)
    sh.WalkSorted()