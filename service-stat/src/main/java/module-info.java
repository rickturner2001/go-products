module goproducts.servicestat {
    requires javafx.controls;
    requires javafx.fxml;

    requires org.kordamp.bootstrapfx.core;
    requires io.github.cdimascio.dotenv.java;
    requires java.logging;
    requires java.sql;

    opens goproducts.servicestat to javafx.fxml;
    exports goproducts.servicestat;
}