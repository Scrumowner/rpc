---
menu:
    after:
        name: graph
        weight: 1
title: Построение графа
---

{{< mermaid >}}
graph LR
    Node1[circle Circle]
    Node2[round-rect Circle]
    Node3[ellipse Circle]
    Node4[circle Rhombus]
    Node5[round-rect Ellipse]
    Node6[ellipse Round Rect]
    Node7[ellipse Square]
    Node8[square Round Rect]
    Node1 --> Node6
    Node2 --> Node1
    Node3 --> Node3
    Node4 --> Node1
    Node5 --> Node8
    Node6 --> Node1
    Node7 --> Node4
    Node8 --> Node5

{{< /mermaid >}}
