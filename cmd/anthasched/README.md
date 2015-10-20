# Draft Design Document for Multi-Element Workflows

DN:

Think of this document as an asynchronous brainstorm session. Everyone is
welcome to edit it. Just check in your comments and changes to this document,
preferably with your initials to help me track the edits.

The basic motivation is to prototype the basic functionality we'll need to
support the execution of multiple antha elements or multiple instances of the
same element. This is a large scope, so a good outcome of this drafting phase
is understanding what will and won't be supported in the near-term, medium-term
or never.

In tension with the desire to minimize the amount of work to achieve our
short-term objectives (i.e., the path to antha v0.5), I'm hoping that with a
few days of thinking about the larger picture, we can come up with a working
model that carries us through to longer term strategic goals (i.e., the path to
antha v2.0).

Initials:
  - DN: Donald Nguyen

Terms:
  - Element: an antha element, the smallest unit of reusable functionality
  - Workflow: a collection of elements into a biological protocol for example
  - Component: a generic term for reusable functionality. An element or a
    workflow, etc.

## Use-cases

A collection of use-cases that we *may* want to address:
  - Clarissa wants to run a construct assembly protocol that produces one
    construct for multiple constructs
    [example](https://github.com/antha-lang/antha/blob/master/antha/component/an/TypeIISConstructAssembly/TypeIISConstructAssembly.an)
  - Elias needs to generate a large number of liquid handling instructions to
    saturate the capabilities of his liquid handling robot.
  - Stephen wants to use an element developed by Dana in his workflow.
  - Victoria wants to validate a workflow developed by David in her own
    lab.
  - Ann wants to analyze the sensitivity of the output of her workflow under
    different parameters (e.g., multi-factorial experiments).
  - Edward wants to make his own workflow but would like to know about existing
    elements that may be useful to him.
  - Peter wants to know what resources would be required to run Siobhan's
    workflow in his lab and what alternative resources or workflows he could
    use to achieve similar outputs.

DN: The most helpful thing for me in the initial brainstorming phase is to
collect a bunch of use-cases. So, if you have any more examples, please pop
them in here!

## Themes / Core Experiences

Besides the Elias example, the use-cases start with an existing workflow and
apply some new capability on top. This suggests to me that *compositionality*
and *discovery* should be core themes.

  - Compositionality: building a new component from existing components without
    knowing too much detail about the constituent components. One implication
    is that components will be used in unforeseen contexts. Classic examples of
    composition are function calls or subtyping in programming languages. A
    more ad-hoc form of composition is copy-and-paste of code. 

  - Discovery: finding components that are responsive to the user and their
    current problem. A tech rule of thumb is that 1% of users will be content
    creators, 9% will be content curators and 90% will be purely content
    consumers. Whether this holds for biology as well, I do not know, but it's
    safe to asssume that the majority of antha users will be using existing
    components with at most a small amount of additional customization to their
    particular needs.

DN: Again, I'd love to hear more input on these topics. The idea is to distill
out a small number of values that characterize a good solution. The main
question is: what are the reoccurring experiences that we hope to cultivate in
order to make users love antha? We want, of course, to validate these guesses
with actual user data, but we need a hypothesis before we can attempt to refute
it.

## Strawman

One possible way to address the use-cases with respect to the core experiences
is, in technical terms, a typed dataflow program or colored Petri network.

Got it?

For those not up to date with programming language theory, dataflow programs
and Petri networks (no relation to the dish) are two graphical ways of modeling
execution. Graphical, in this case, refers to the underlying mathematic
abstraction, a graph. A graph is a set of nodes (think circles) and a set of
edges (think arrows) that connect nodes together. In both dataflow and Petri
nets, nodes represent computations and edges represent the flow of information
from one computation to another. For example, an add node may take two inputs
(i.e., its summands), represented as edges pointing to the node, and the add
node produces a sum, represented as one edge pointing away from the node.

In both dataflow and Petri nets, a particular instance of a computation (like
computing the sum of 3 and 5) is modeled by tokens that flow through the graph,
and a computation can occur when tokens arrive on all the incoming edges of a
node.  Thus, tokens representing the summands flow into the node, the sum node
performs its computation, and the result flows out from the node. The main
difference between dataflow programs and Petri nets is where these tokens sit.
In dataflow programs, tokens sit on the edges between nodes in a queue, while
in Petri nets, tokens sit on nodes themselves.

This seemingly small difference actually has a huge impact in the behavior of
the two systems. In Petri nets, there is no particular order in which tokens
are accumulated on nodes, so if there are multiple downstream computations that
could consume a token, which one fires is a choice of the system implementing
the Petri network. Petri networks are inherently non-deterministic. In
dataflow, tokens line up in orderly queues on the edges between nodes, and it
turns out that regardless of the order in which computations consume their
tokens, the overall output of a dataflow program will always be the same.

This graphical way of modeling computation occurs all over the place under
various names, dependence graphs, Gantt charts, graphical models, enzyme
pathway charts, etc. So, it appears to be a very natural way to represent
computation. But as anyone who as looked at a graphical model of any modest
size can attest, they can quickly become unwieldy and hard to understand.

There are a couple of things we can do to manage the complexity. 

  1. Abstraction. Any connected subgraph in a graphical model can be abstracted
     by collapsing the subgraph into a single node. This imposes a hierarchical
     structure on the graph.
     
     "Antha give me a bird's eye view of this workflow. Ok, now zoom in on
     this area."

  2. Typing. We can give tokens and edges colors such that tokens only flow
     along matching edges, and we can give edges names. Given a set of
     components, there will be fewer ways (hopefully only one way) in which
     components can be assembled, which means the tedious task of wiring up
     components can be automated. For example, there is only one way of wiring
     up a construct assembly, transformation and screening protocol. 

     The fact that a component may be identifiable based on its graphical
     structure alone (i.e., the number and types of its edges) can serve as a
     basis for component discovery.
     
     "Antha find me any component that can fill this hole in my workflow."

Typed dataflow and colored Petri networks are well-known general-purpose
abstractions. To this, I propose extending the model in a few ways to better
address some of the use-cases. 

  1. Characterization of the domain of external inputs. Graphically, parameters
     and inputs external to a component are represented as nodes with no
     incoming edges. They are just token emitters. Some external inputs are
     under control of the system (e.g., parameters to the workflow) and some
     are under control of the environment (e.g., sensor readings). Since
     protocols are usually dependent on the particular values these parameters
     may take (e.g., a protocol that works on 1L of liquid is different than
     one that works on 1ul of liquid), it makes sense for components to
     characterize what range of values it expects for each external parameter
     (e.g., a range 0-30C or a distribution 5ul +- 1%). This information will
     be used in component discovery, workflow simulation and might be used for
     simple sanity checking of workflows.

  2. Parameter exploration policies. Instead of emitting a particular stream of
     tokens, external inputs emit non-deterministic samples of tokens. This
     decouples a workflow definition from a particular execution (i.e., actual
     tokens consumed). Particular executions can be specified on the side
     either at a high-level (e.g., "Antha, give me an execution where parameter
     X varies the most frequently.") or low-level.

  3. Variable token consumption. This is more of a technical detail, but in
     order to more precisely model components, a node may consume a non-unit
     number of tokens on its input edges. This helps model loops in the graph
     and loops in element programs. If the number of tokens consumed is not
     bounded by the component, it should be a value chosen by the antha system.

### Composition Strategies

  - Parameter exploration. "Antha, run this workflow at 25C." Since parameter
    values are not baked into workflows, it should be easy to explore a few
    extra points in the parameter space.
  
  - Graph assembly. "Antha, please connect these two components." Since there
    are only a few ways two components can be wired together, we reduce the
    tedious act of wiring with selection of a few alternatives (hopefully).

  - Graph substitution. "Antha, use this workflow but replace subgraph A with
    subgraph B." This style of composition can be very verbose because
    descriptions of arbitary subgraphs are the same size as the subgraphs
    themselves. Hopefully, we can identify a few concise patterns. Perhaps
    typing could help here. There may only be a few places to insert a
    componenet. "Antha, use this workflow with this component instead." 

DN: What other strategies would be common?

### Discovery Strategies

DN: Not sure if my assumption that biological protocols can be distinguished 
based on their types is a realistic one.

### Execution Model

The programming model is in terms of logical operations in a virtual
(unbounded) space. To actually execute an antha program, we need to lower it
into a representation on physical (bounded) resources and bulk up individual
equipment instructions into plate-wise operations.  Equipment can be modeled as
an external input that emits a stream of tickets that enable the consumer to
perform one instance of a desired action. For now, let's use a simple directed
acyclic graph as the lowered representation. DAG nodes are actions labeled with
physical resources (e.g., a plate) and token values from the logical execution
model.

Lowering is iterative application of, for lack of a better term, a process I'll
call pumping. Pumping is calculating what tokens would be needed to produce a
desired output token. Imagine the suction effect one gets when pumping water.
The idea is to pump enough output tokens to get enough tokens from a piece of
equipment (these will be external inputs that must be ancestors of the node
producing the output token) to form a block operation.  If all antha operations
had inverses, we could calculate this directly, but since they likely don't
(examples of things that cannot be modeled precisely: loops, components that
produce tokens on different outputs based on input values), pumping computes an
approximation of the required tokens by following the graph dependencies,
applying some static/symbolic analysis (TBD) and potentially "bubble"
instructions/resources which are noop instructions or unused resources that
need to be introduced because pumping could not determine the input
dependencies precisely enough.

## Scratch Pad

Unstructured notes go here.

Error handling / Validation errors
