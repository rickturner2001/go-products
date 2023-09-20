package goproducts.servicestat;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.logging.Level;
import java.util.logging.Logger;

public class ServicesDatabase {
    Connection connection;
    Statement statement;
    Logger logger;


    public ServicesDatabase() {

        this.logger = Logger.getLogger(Manager.class.getName());
        logger.setLevel(Level.FINE);


        try {

            this.connection = DriverManager.getConnection("jdbc:sqlite:services.db");
            this.statement = connection.createStatement();
            statement.setQueryTimeout(15);  // set timeout to 30 sec.
            statement.executeUpdate("create table if not exists Service(name varchar, url varchar)");
        } catch (SQLException err) {
            logger.log(Level.SEVERE, "could not execute creation statement", err);
        }
    }

    public boolean addService(Service service) {

        String statementString = String.format("insert into Service values('%s', '%s')", service.getName(), service.getUrlString());
        try {
            statement.executeUpdate(statementString);
            return true;
        } catch (SQLException err) {
            this.logger.log(Level.SEVERE, String.format("could not execute statement %s", statementString));
            return false;
        }

    }

    public ArrayList<Service> getServices() {

        ArrayList<Service> services = new ArrayList<>();
        try {
            ResultSet rs = statement.executeQuery("select * from Service");

            while (rs.next()) {
                String name = rs.getString("name");
                String url = rs.getString("url");

                services.add(new Service(name, url));
            }
        } catch (
                SQLException err) {
            this.logger.log(Level.SEVERE, "could not get services from db", err);
        }

        return services;
    }
}
