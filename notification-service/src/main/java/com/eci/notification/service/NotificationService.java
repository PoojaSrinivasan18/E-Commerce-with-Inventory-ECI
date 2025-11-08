package com.eci.notification.service;

import com.eci.notification.model.Notification;
import com.eci.notification.repository.NotificationRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.LocalDateTime;

@Service
public class NotificationService {

	@Autowired
	private NotificationRepository notificationRepository;

	@Autowired
	private JavaMailSender mailSender;
	
	private static final Logger logger = LoggerFactory.getLogger(NotificationService.class);

	// --- Explicit DB call method ---
	public Notification logNotificationToDB(Notification notification) {
		return notificationRepository.save(notification);
	}

	public Notification sendEmail(String toEmail, String subject, String body, Long orderId, Long paymentId,
			Long shipmentId, String type, String channel, String messageContent) {
		Notification notification = new Notification();
		notification.setOrderId(orderId);
		notification.setPaymentId(paymentId);
		notification.setShipmentId(shipmentId);
		notification.setType(type != null ? type : "ORDER");
		notification.setChannel(channel != null ? channel : "EMAIL");
		notification.setMessageContent(messageContent != null ? messageContent : body);

		try {
			logger.info("Attempting to send email to: {} | subject: {} | body: {}", toEmail, subject, body);
			SimpleMailMessage message = new SimpleMailMessage();
			message.setTo(toEmail);
			message.setSubject(subject);
			message.setText(body);
			mailSender.send(message);

			notification.setStatus("SENT");
			notification.setSentAt(LocalDateTime.now());
			logger.info("Email sent successfully to: {}", toEmail);
		} catch (Exception e) {
			notification.setStatus("FAILED");
			 logger.error("Failed to send email to: {} | error: {}", toEmail, e.getMessage(), e);
		}
		
		logger.info("Saving notification to DB: {}", notification);
		// --- Save to DB ---
		return logNotificationToDB(notification);
	}

}
