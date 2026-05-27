package com.shopee.auth.crypto;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class PiiAttributeConverter implements AttributeConverter<String, String> {

    @Override
    public String convertToDatabaseColumn(String attribute) {
        if (attribute == null || attribute.isEmpty()) {
            return attribute;
        }
        EncryptionUtil encryption = EncryptionHolder.getInstance();
        if (encryption != null && encryption.isEnabled()) {
            return encryption.encrypt(attribute);
        }
        return attribute;
    }

    @Override
    public String convertToEntityAttribute(String dbData) {
        if (dbData == null || dbData.isEmpty()) {
            return dbData;
        }
        EncryptionUtil encryption = EncryptionHolder.getInstance();
        if (encryption != null && encryption.isEnabled()) {
            return encryption.decrypt(dbData);
        }
        return dbData;
    }
}
