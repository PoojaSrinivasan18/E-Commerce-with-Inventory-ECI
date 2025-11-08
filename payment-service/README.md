# Payment Service
The Payment Service is a dedicated microservice responsible for handling all payment-related operations within the E-commerce with Inventory (ECI) platform. It follows a database-per-service pattern, managing its own isolated Payment database, which stores payment transactions and idempotency keys to ensure consistency and reliability.

# Key Responsibilities
* Charge and refund payments associated with customer orders, ensuring idempotency on all charge and refund operations to prevent duplicate transactions.
* Update payment status on orders, triggering downstream workflows such as order confirmation or cancellation handling.
* Handle idempotent charge requests using unique idempotency keys, which guarantee safe retries without double charging.
* Support both immediate charge mode and potential extension to authorize-capture flows (for advanced fulfillment scenarios).
* Provide APIs to initiate charges, refunds, and query payment status.
* Ensure transactional integrity and robust error handling, including retry policies with jitter for network or transient failures.
* Emit events or update order statuses upon successful or failed payments to coordinate with Order and Inventory services.
* Expose a versioned REST API /v1/payments following OpenAPI 3.0 standards with comprehensive error schemas, pagination, and filtering.

# Database Schema (Database-Per-Service)
* Payments: Stores payment records with fields such as payment_id, order_id, amount, method, status, reference, and created_at.
* Idempotency Keys: Tracks unique keys for idempotent operations to avoid duplicate processing.

