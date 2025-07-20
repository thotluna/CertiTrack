# Diseño de la Base de Datos - Sistema CertiTrack

## 1. Introducción

Este documento detalla el diseño propuesto para la base de datos del Sistema de Gestión de Certificaciones (CertiTrack). El objetivo principal es establecer una estructura de datos robusta, escalable y flexible que soporte los requisitos funcionales identificados, permitiendo la gestión eficiente de certificaciones asociadas a personas y equipos, así como el manejo de información adicional variable para cada certificación y el envío flexible de notificaciones.

Se ha buscado un equilibrio entre la normalización relacional y la flexibilidad para adaptarse a la diversidad de tipos de certificaciones y sus atributos específicos.

## 2. Decisiones Clave de Diseño

### 2.1. Gestión de Registros (Soft Delete)

**Decisión:** Se implementará la "eliminación lógica" (soft delete) para los registros de personas, equipos y certificaciones, en lugar de la eliminación física de los datos.

**Justificación:**
* **Integridad Referencial y Historial:** Permite mantener el historial completo de certificaciones y sus asociaciones, incluso si una persona egresa o un equipo es dado de baja. Esto es crucial para auditorías y para comprender el contexto de certificaciones pasadas.
* **Recuperación de Datos:** Facilita la recuperación de registros eliminados accidentalmente.
* **Cumplimiento y Auditoría:** Soporta requisitos legales o de cumplimiento que exigen la retención de datos por períodos específicos.
* **Evita Inconsistencias:** Previene la ruptura de relaciones de datos si un registro es eliminado.

**Implementación:** Se añadirá una columna de `estado` (o similar, ej., `activo`, `fecha_baja`) a las tablas correspondientes (`Persona`, `Equipo`, `Certificacion`). La "eliminación" de un registro implicará cambiar este estado (ej., de `ACTIVO` a `INACTIVO` o `DADO_DE_BAJA`), sin eliminar la fila de la base de datos.

### 2.2. Manejo de Asociaciones de Certificaciones (Asociación Polimórfica)

**Problema Identificado:** Las certificaciones pueden estar asociadas a una `Persona`, a un `Equipo`, o incluso en el futuro a otro tipo de entidad. Intentar usar `id_persona` y `id_equipo` directamente en la tabla `Certificacion` crearía columnas nulas innecesarias y requeriría una lógica de aplicación compleja para asegurar que solo una de las dos FKs estuviera presente, o al menos una.

**Decisión:** Se implementará un patrón de "asociación polimórfica" para vincular una certificación a su entidad principal (Persona o Equipo).

**Justificación:**
* **Flexibilidad:** Permite que una certificación se asocie a diferentes tipos de entidades (Persona o Equipo) de manera limpia.
* **Extensibilidad:** Facilita la adición de nuevos tipos de entidades asociadas en el futuro sin modificar la tabla principal `Certificacion`.
* **Evita Columnas Nulas:** Reduce la cantidad de columnas `NULL` en la tabla `Certificacion`.

**Implementación:** La tabla `Certificacion` contendrá dos nuevas columnas:
* `entidad_asociada_id`: Almacenará el identificador único (ID) de la persona o el equipo al que la certificación está asociada. Este campo no tendrá una Foreign Key directa a `Persona` o `Equipo` a nivel de base de datos.
* `tipo_entidad_asociada`: Almacenará un valor que indicará el tipo de entidad a la que `entidad_asociada_id` hace referencia (ej., 'PERSONA', 'EQUIPO').

La **lógica de la aplicación** será responsable de:
* Validar que `entidad_asociada_id` contenga un ID válido para el `tipo_entidad_asociada` especificado.
* Asegurar que se cumplan las reglas de negocio (ej., que una certificación de alimentos solo se asocie a una `PERSONA`, o que una de calibración solo a un `EQUIPO`).

### 2.3. Manejo de Datos Adicionales/Atributos Específicos de Certificación

**Problema Identificado:** Diferentes tipos de certificaciones requieren almacenar datos específicos y variables (ej., "Institución Emisora", "Número de Horas", "Rango de Calibración"). Añadir una columna por cada posible atributo a la tabla `Certificacion` resultaría en una tabla muy ancha, con muchas columnas nulas y sería rígida para futuros cambios.

**Decisión:** Se creará una tabla `CertificacionAtributo` para almacenar estos datos adicionales de forma flexible utilizando un esquema clave-valor.

**Justificación:**
* **Flexibilidad Extrema:** Permite asociar cualquier número de atributos personalizados a cualquier certificación sin modificar el esquema de la tabla `Certificacion` principal.
* **Normalización:** Mejora la estructura de la base de datos al evitar columnas vacías y agrupar datos relacionados de manera eficiente.
* **Extensibilidad:** Facilita la incorporación de nuevos tipos de certificaciones con sus propios atributos específicos sin necesidad de cambios en el esquema de la base de datos.

**Implementación:** Se añadirá una nueva tabla llamada `CertificacionAtributo`.

### 2.4. Gestión de Múltiples Destinatarios para Notificaciones

**Problema Identificado:** Las notificaciones pueden requerir ser enviadas a múltiples destinatarios (ej., al titular de la certificación, al jefe de la unidad, a un contacto de equipo), y registrar el estado de envío para cada uno.

**Decisión:** Se creará una tabla `AlertaDestinatario` para almacenar los detalles de cada destinatario de una notificación específica, permitiendo una relación uno a muchos.

**Justificación:**
* **Flexibilidad de Destinatarios:** Permite asociar múltiples direcciones de correo electrónico a una única instancia de notificación.
* **Trazabilidad Individual:** Permite registrar y monitorear el estado de envío y la fecha para cada destinatario de forma independiente.
* **Claridad del Modelo:** Separa la información de la notificación general de los detalles de sus destinatarios específicos.

**Implementación:**
* El campo `destinatario_correo` se eliminará de la tabla `AlertaNotificacion`.
* Se añadirá una nueva tabla `AlertaDestinatario` que contendrá el `id_alerta_notificacion` (FK), el `correo_electronico`, el `tipo_destinatario`, `estado_envio_individual` y `fecha_envio_individual`.

## 3. Estructura de Entidades y Atributos (Modelo Lógico)

A continuación, se describen las entidades principales de la base de datos con sus atributos y tipos de datos propuestos (genéricos, a adaptar a la base de datos específica como MySQL, PostgreSQL, SQL Server, etc.).

### 3.1. Entidad: Persona

Representa a los individuos que pueden poseer certificaciones.

* `id`: Identificador único de la persona (VARCHAR(50), PK).
* `nombre`: Nombre de la persona (VARCHAR(100)).
* `apellido`: Apellido de la persona (VARCHAR(100)).
* `identificacion`: Cédula o documento de identidad (VARCHAR(20), Único).
* `correo_electronico`: Correo electrónico de contacto (VARCHAR(255)).
* `telefono`: Número de teléfono (VARCHAR(20)).
* `estado`: Estado actual de la persona (VARCHAR(20), ej., 'ACTIVO', 'INACTIVO', 'EGRESADO').
* `fecha_creacion`: Marca de tiempo de la creación del registro (DATETIME).
* `fecha_actualizacion`: Última marca de tiempo de actualización del registro (DATETIME).

### 3.2. Entidad: Equipo

Representa los equipos o activos que pueden tener certificaciones asociadas.

* `id`: Identificador único del equipo (VARCHAR(50), PK).
* `nombre`: Nombre o modelo del equipo (VARCHAR(100)).
* `descripcion`: Descripción detallada del equipo (TEXT).
* `numero_activo`: Número de inventario o activo fijo (VARCHAR(50), Opcional, Único).
* `estado`: Estado actual del equipo (VARCHAR(20), ej., 'ACTIVO', 'DADO_DE_BAJA', 'EN_MANTENIMIENTO').
* `fecha_creacion`: Marca de tiempo de la creación del registro (DATETIME).
* `fecha_actualizacion`: Última marca de tiempo de actualización del registro (DATETIME).
* `id_responsable_persona`: Clave Foránea (FK) a `Persona`. Identifica a la persona principal responsable del equipo. (VARCHAR(50), Opcional).

### 3.3. Entidad: TipoCertificacion

Define las categorías o tipos predefinidos de certificaciones que el sistema puede manejar.

* `id`: Identificador único del tipo de certificación (VARCHAR(50), PK).
* `nombre`: Nombre del tipo de certificación (VARCHAR(100), ej., 'AWS Certified Cloud Practitioner', 'Certificado Manipulación Alimentos').
* `descripcion`: Descripción del tipo de certificación (TEXT).

### 3.4. Entidad: Certificacion

La entidad central que representa cada certificación registrada en el sistema.

* `id`: Identificador único de la certificación (VARCHAR(50), PK).
* `entidad_asociada_id`: ID de la entidad a la que está asociada la certificación (Persona o Equipo). (VARCHAR(50)).
* `tipo_entidad_asociada`: Indica el tipo de entidad asociada ('PERSONA', 'EQUIPO'). (VARCHAR(20)).
* `id_tipo_certificacion`: Clave Foránea (FK) a `TipoCertificacion`. (VARCHAR(50)).
* `nombre_certificacion`: Nombre específico de la certificación si difiere del tipo (VARCHAR(255)).
* `codigo_certificacion`: Código o número de referencia de la certificación (VARCHAR(100), Opcional).
* `fecha_emision`: Fecha en que la certificación fue emitida (DATE).
* `fecha_vencimiento`: Fecha en que la certificación expira (DATE, Opcional, si la certificación no vence).
* `estado`: Estado actual de la certificación (VARCHAR(20), ej., 'ACTIVA', 'VENCIDA', 'CANCELADA', 'REVOCADA').
* `fecha_creacion`: Marca de tiempo de la creación del registro (DATETIME).
* `fecha_actualizacion`: Última marca de tiempo de actualización del registro (DATETIME).

### 3.5. Entidad: Documento

Almacena la metadata de los documentos adjuntos a las certificaciones. Los archivos físicos no se guardan en la DB, sino en un sistema de almacenamiento externo.

* `id`: Identificador único del documento (VARCHAR(50), PK).
* `id_certificacion`: Clave Foránea (FK) a `Certificacion`. (VARCHAR(50)).
* `nombre_archivo`: Nombre original del archivo (VARCHAR(255)).
* `ruta_almacenamiento`: Ruta o URL donde el archivo está almacenado en el sistema externo (VARCHAR(512)).
* `tipo_mime`: Tipo MIME del archivo (ej., 'application/pdf', 'image/jpeg') (VARCHAR(50)).
* `tamano_bytes`: Tamaño del archivo en bytes (BIGINT).
* `fecha_subida`: Marca de tiempo de la subida del documento (DATETIME).

### 3.6. Entidad: AlertaNotificacion

Representa una instancia de alerta programada, sin especificar aún sus destinatarios individuales.

* `id`: Identificador único de la alerta (VARCHAR(50), PK).
* `id_certificacion`: Clave Foránea (FK) a `Certificacion`. (VARCHAR(50)).
* `fecha_programada`: Fecha y hora en que la alerta fue programada para ser enviada (DATETIME).
* `fecha_envio_general`: Fecha y hora real de inicio del proceso de envío general de la alerta (DATETIME, NULL si no se ha iniciado/completado).
* `tipo_alerta`: Tipo de alerta (VARCHAR(50), ej., 'VENCIMIENTO_PROXIMO', 'VENCIMIENTO_HOY').
* `estado_general`: Estado general de la notificación (VARCHAR(20), ej., 'PENDIENTE', 'ENVIADA_PARCIAL', 'ENVIADA_COMPLETA', 'FALLIDA').

### 3.7. Entidad: AlertaDestinatario

Almacena los detalles de cada correo electrónico al que se debe enviar o se envió una instancia específica de `AlertaNotificacion`.

* `id`: Identificador único del destinatario de alerta (VARCHAR(50), PK).
* `id_alerta_notificacion`: Clave Foránea (FK) a `AlertaNotificacion`. (VARCHAR(50)).
* `correo_electronico`: Dirección de correo electrónico específica para este envío (VARCHAR(255)).
* `tipo_destinatario`: Rol o tipo de este destinatario (VARCHAR(50), ej., 'PRINCIPAL', 'COPIA', 'JEFE_UNIDAD', 'RESPONSABLE_EQUIPO', 'CONTACTO_ADICIONAL').
* `estado_envio_individual`: Estado del envío para este destinatario particular (VARCHAR(20),