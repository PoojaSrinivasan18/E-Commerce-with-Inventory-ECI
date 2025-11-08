# Catalog Service
The Catalog Service is a dedicated microservice responsible for managing all product-related operations within the E-commerce with Inventory (ECI) platform. It follows a database-per-service architecture, maintaining its own isolated Catalog database that stores detailed information about products, including their pricing, category, and availability status. This ensures data autonomy, scalability, and flexibility for independent service evolution.

# Key Responsibilities

* Maintain and manage product information, including SKU, name, category, price, and active status.

* Support full CRUD operations (Create, Read, Update, Delete) for product management.

* Provide efficient search, filtering, and pagination capabilities to fetch products based on criteria like category, name, and price range.

* Ensure product data consistency while allowing replication or synchronization with other services such as Inventory or Order when required.

* Handle product availability and pricing queries through lightweight, optimized APIs.

* Expose a versioned REST API /v1/catalog following OpenAPI 3.0 standards with consistent error schemas, structured responses, and metadata.

* Facilitate downstream services (Inventory, Orders, Payments) to reference products through secure and stable APIs.

* Optionally publish product update events (product.created, product.updated, product.deleted) for other services to stay synchronized.

# Database Schema (Database-Per-Service)

Products Table: Stores all product-related details.
Fields include:

* product_id – Unique identifier for each product.

* sku – Unique product code for external reference and consistency.

* name – Product name or title.

* category – Category under which the product is classified.

* price – Product price, stored with appropriate precision.

* is_active – Boolean flag indicating product availability for sale.

* description – Additional details about the product.

* created_at – Timestamp of when the product was added.

* updated_at – Timestamp of the latest update to the product details.