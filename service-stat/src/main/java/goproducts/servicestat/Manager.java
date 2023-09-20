package goproducts.servicestat;

import io.github.cdimascio.dotenv.Dotenv;

import java.io.IOException;
import java.lang.reflect.Array;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Objects;
import java.util.logging.Level;
import java.util.logging.Logger;


public class Manager {


    private final Logger logger;
    private ArrayList<Service> services;


    public Manager(ArrayList<Service> services) {
        this.logger = Logger.getLogger(Manager.class.getName());
        logger.setLevel(Level.FINE);

        this.services = services;


        this.logger.log(Level.INFO, "Successfully initialized Manager");
    }


    public static void main(String[] args) {

        Dotenv dotenv = Dotenv.load();
        ArrayList<Service> services = new ArrayList<>();

        services.add(new Service("product-api", "http://localhost:8080/status"));
        services.add(new Service("postgres-db", dotenv.get("DATABASE_URL")));

        Manager manager = new Manager(services);

        manager.checkServicesStatus();

        int activeServices = manager.activeServices();

        manager.logger.log(Level.INFO, String.format("There are %d active services\n", activeServices));
    }


    public void setServices(ArrayList<Service> services) {
        this.services = services;
    }

    public Service addService(Service service) {

        this.services.add(service);
        return service;
    }

    public ArrayList<Service> getServices() {
        return services;
    }

    public void checkServicesStatus() {
        for (Service service : this.services) {
            checkServiceStatus(service);
        }
    }


    public int activeServices() {
        int activeServicesCounter = 0;
        for (Service service : this.services) {
            if (service.getStatus() == ServiceStatus.ACTIVE) {
                activeServicesCounter++;
            }
        }
        return activeServicesCounter;
    }

    private void checkServiceStatus(Service service) {
        boolean isActive = sendRequest(service);

        if (isActive) {
            this.logger.log(Level.INFO, String.format("Service %s is active\n", service.getName()));
            service.setStatusActive();
        } else {

            this.logger.log(Level.INFO, String.format("Service %s is inactive\n", service.getName()));
            service.setStatusInactive();
        }

    }

    private boolean sendRequest(Service service) {
        try {
            URL url = new URL(service.urlString);
            HttpURLConnection con = (HttpURLConnection) url.openConnection();

            con.setRequestMethod("GET");

            int responseCode = con.getResponseCode();

            return responseCode == HttpURLConnection.HTTP_OK;
        } catch (MalformedURLException e) {
            this.logger.log(Level.SEVERE, String.format("Could not parse URL: %s\n%s\n", service.urlString, e));
            service.setStatusInactive();
            return false;
        } catch (IOException e) {
            this.logger.log(Level.SEVERE, String.format("Could instantiate openConnection: %s\n", e));
            service.setStatusInactive();
            return false;
        }


    }
}

