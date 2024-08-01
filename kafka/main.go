// Kafka service module for Dagger.
package main

import (
	"dagger/kafka/internal/dagger"
	"fmt"
)

// defaultImageRepository is used when no image is specified.
const defaultImageRepository = "bitnami/kafka"

type Kafka struct {
	Container *dagger.Container
}

func New(
	// Version (image tag) to use from the official image repository as a base container.
	//
	// +optional
	version string,

	// Custom container to use as a base container.
	//
	// +optional
	container *dagger.Container,
) *Kafka {
	if container == nil {
		if version == "" {
			version = "latest"
		}

		container = dag.Container().From(fmt.Sprintf("%s:%s", defaultImageRepository, version))
	}

	container = container.
		// https://github.com/bitnami/charts/issues/22552#issuecomment-1905721850
		WithEnvVariable("KAFKA_CFG_MESSAGE_MAX_BYTES", "1048588").

		// KRaft settings
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@127.0.0.1:9093").
		// Listeners
		WithEnvVariable("KAFKA_CFG_LISTENERS", "PLAINTEXT://:9092,CONTROLLER://:9093").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "PLAINTEXT://kafka:9092").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "PLAINTEXT").
		WithExposedPort(9092)

	return &Kafka{container}
}

// Set an environment variable.
func (m *Kafka) WithEnvVariable(
	// The name of the environment variable (e.g., "HOST").
	name string,

	// The value of the environment variable (e.g., "localhost").
	value string,

	// Replace `${VAR}` or $VAR in the value according to the current environment
	// variables defined in the container (e.g., "/opt/bin:$PATH").
	//
	// +optional
	expand bool,
) *Kafka {
	return &Kafka{
		m.Container.WithEnvVariable(name, value, dagger.ContainerWithEnvVariableOpts{
			Expand: expand,
		}),
	}
}

// Launch a Kafka service.
func (m *Kafka) Service() *dagger.Service {
	return m.Container.AsService()
}
