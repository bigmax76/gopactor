/*
Package Gopactor provides a set of tools to simplify testing
of actors created with Protoactor (https://github.com/AsynkronIT/protoactor-go).

Main features:

Intercept messages

For any actor you want to test, Gopactor can intercept all it's inbound and outbound
messages. It is probably exactly what you want to do when you test the actor's behavior.
Moreover, interception forces a naturally asynchronous actor to act in a more synchronous way.
When messages are sent and received under the control of Gopactor, it is much easier to reason about
the actor's logic and examine its communication with the outside world step by step.

Intercept system messages

Protoactor uses special system messages to control the lifecycle of an actor.
Gopactor can intercept some of such messages to help you ensure that your actor
stops or restarts when expected.

Intercept spawning of children

It is a common pattern to let actors spawn child actors and communicate with them.
Good as it is, this pattern often stays in the way of writing deterministic tests.
Given that child-spawning and communication happen in the background asynchronously,
it can be seen more like a side-effect that can interfere with your tests in many
unpredictable ways.

By default, Gopactor intercepts all spawn invocations and instead of spawning what is
requested, it spawns no-op null-actors. These actors are guaranteed to not communicate
with their parents in any way. If you do no want Gopactor to substitute spawned actors,
you can easily disable this behavior via configuration options.

Goconvey-style assertions

Gopactor provides a bunch of assertion functions to be used with the popular testing
framework Goconvey (http://goconvey.co/). For instance,

	So(worker, ShouldReceive, "ping")
	So(worker, ShouldSendTo, requestor, "pong")

Configurable

For every tested actor, you can define what you want to intercept: inbound, outbound
or system messages. Or spawning of children. Or everything. Or nothing at all.
You can also set a custom timeout:

	options := OptNoInterception.
		WithOutboundInterception().
		WithPrefix("my-actor").
		WithTimeout(10 * time.Millisecond)

Example of usage

Here is a short example. We'll define and test a simple worker actor that can do only one thing:
respond "pong" when it receives "ping".

	package worker_test

	import (
		"testing"

		"github.com/AsynkronIT/protoactor-go/actor"
		. "github.com/meAmidos/gopactor"
		. "github.com/smartystreets/goconvey/convey"
	)

	// Actor to test
	type Worker struct{}

	// This actor is very simple. It can do only one thing, but it does this thing well.
	func (w *Worker) Receive(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case string:
			if m == "ping" {
				ctx.Respond("pong")
			}
		}
	}

	func TestWorker(t *testing.T) {
		Convey("Test the worker actor", t, func() {

			// It is essential to spawn the tested actor using Gopactor. This way, Gopactor
			// will be able to intercept all inbound/outbound messages of the actor.
			worker, err := SpawnFromInstance(&Worker{}, OptDefault.WithPrefix("worker"))
			So(err, ShouldBeNil)

			// Spawn an additional actor that will communicate with our worker.
			// The only purpose of this actor is to be a sparring partner,
			// so we don't care about its functionality.
			// Conveniently, Gopactor provides an easy way to create it.
			requestor, err := SpawnNullActor()
			So(err, ShouldBeNil)

			// Let the requestor ping the worker
			worker.Request("ping", requestor)

			// Assert that the worker receives the ping message
			So(worker, ShouldReceive, "ping")

			// Assert that the worker sends back the correct response
			So(worker, ShouldSendTo, requestor, "pong")

			// Finally, assert that the requestor gets the response
			So(requestor, ShouldReceive, "pong")
		})
	}
*/
package gopactor
