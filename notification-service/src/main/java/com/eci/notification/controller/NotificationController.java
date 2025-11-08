package com.eci.notification.controller;

import com.eci.notification.model.Notification;
import com.eci.notification.service.NotificationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@RestController
@RequestMapping("/v1/notifications")
public class NotificationController {
	private static final Logger logger = LoggerFactory.getLogger(NotificationController.class);

	@Autowired
	private NotificationService notificationService;

	@PostMapping("/email")
	public Notification sendEmail(@RequestBody EmailRequest request) {
		logger.info("Received request to /v1/notifications/email with data: {}", request);
		return notificationService.sendEmail(
				request.getCustomerEmail(),
				request.getSubject(),
				request.getMessage(),
				request.getOrderId(),
				request.getPaymentId(),
				request.getShipmentId(),
				request.getType(),
				request.getChannel(),
				request.getMessageContent());
	}
}

// --- DTO Classes with Getters & Setters ---

class EmailRequest {
	private Long orderId;
	private Long paymentId;
	private Long shipmentId;
	private String type; // ORDER, PAYMENT, SHIPMENT
	private String channel; // EMAIL
	private String status; // SENT, FAILED, PENDING
	private String messageContent;
	private String customerEmail;
	private String subject;
	private String message;

	public Long getOrderId() {
		return orderId;
	}

	public void setOrderId(Long orderId) {
		this.orderId = orderId;
	}

	public Long getPaymentId() {
		return paymentId;
	}

	public void setPaymentId(Long paymentId) {
		this.paymentId = paymentId;
	}

	public Long getShipmentId() {
		return shipmentId;
	}

	public void setShipmentId(Long shipmentId) {
		this.shipmentId = shipmentId;
	}

	public String getType() {
		return type;
	}

	public void setType(String type) {
		this.type = type;
	}

	public String getChannel() {
		return channel;
	}

	public void setChannel(String channel) {
		this.channel = channel;
	}

	public String getStatus() {
		return status;
	}

	public void setStatus(String status) {
		this.status = status;
	}

	public String getMessageContent() {
		return messageContent;
	}

	public void setMessageContent(String messageContent) {
		this.messageContent = messageContent;
	}

	public String getCustomerEmail() {
		return customerEmail;
	}

	public void setCustomerEmail(String customerEmail) {
		this.customerEmail = customerEmail;
	}

	public String getSubject() {
		return subject;
	}

	public void setSubject(String subject) {
		this.subject = subject;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	@Override
	public String toString() {
		return "EmailRequest{" + "orderId=" + orderId + ", paymentId=" + paymentId + ", shipmentId=" + shipmentId
				+ ", type='" + type + '\'' + ", channel='" + channel + '\'' + ", status='" + status + '\''
				+ ", messageContent='" + messageContent + '\'' + ", customerEmail='" + customerEmail + '\''
				+ ", subject='" + subject + '\'' + ", message='" + message + '\'' + '}';
	}
}