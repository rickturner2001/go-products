package goproducts.servicestat;

import io.github.cdimascio.dotenv.Dotenv;
import javafx.collections.FXCollections;
import javafx.collections.ListChangeListener;
import javafx.collections.ObservableList;
import javafx.event.ActionEvent;
import javafx.fxml.FXML;
import javafx.scene.control.TableColumn;
import javafx.scene.control.TableView;
import javafx.scene.control.TextField;
import javafx.scene.control.cell.PropertyValueFactory;
import javafx.scene.layout.HBox;
import javafx.scene.layout.VBox;
import javafx.scene.text.Text;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

public class Controller {
    public TextField serviceName;
    public TextField serviceUrl;
    private ObservableList<Service> activeServices;
    private final Manager manager;
    public HBox alertContainer;
    @FXML
    public VBox tableViewContainer;

    private ServicesDatabase database;


    public Controller() {


        Dotenv dotenv = Dotenv.load();
        this.database = new ServicesDatabase();

        ArrayList<Service> services = new ArrayList<>();
        this.manager = new Manager(services);


    }

    @FXML
    private void initialize() {
        TableView<Service> serviceTable = new TableView<>();

        serviceTable.setPrefWidth(460);

        TableColumn<Service, Object> serviceNameColumn = new TableColumn<>("Name");
        serviceNameColumn.setCellValueFactory(new PropertyValueFactory<>("name"));

        serviceNameColumn.setPrefWidth((double) 460 / 3);

        TableColumn<Service, Object> serviceUrlColumn = new TableColumn<>("URL");
        serviceUrlColumn.setCellValueFactory(new PropertyValueFactory<>("urlString"));
        serviceUrlColumn.setPrefWidth((double) 460 / 3);

        TableColumn<Service, Object> serviceStatusColumn = new TableColumn<>("Status");
        serviceStatusColumn.setCellValueFactory(new PropertyValueFactory<>("status"));
        serviceStatusColumn.setPrefWidth((double) 460 / 3);

        serviceTable.getColumns().add(serviceNameColumn);
        serviceTable.getColumns().add(serviceUrlColumn);
        serviceTable.getColumns().add(serviceStatusColumn);

        this.activeServices = FXCollections.observableArrayList();


        tableViewContainer.getChildren().add(serviceTable);

        var dbServices = database.getServices();
        manager.setServices(dbServices);


        for (Service service : manager.getServices()) {
            System.out.printf("Service: %s%n", service.getName());
            activeServices.add(service);
        }

        serviceTable.setItems(activeServices);


    }

    public void addServiceToObservable(Service inputService) {


        for (Service service : this.manager.getServices()) {
            if (service.getName().equals(inputService.getName())) {
                System.out.println(String.format("Service with name %s already exists", inputService.getName()));
                return;
            }
        }


        this.activeServices.add(inputService);
        this.manager.addService(inputService);
        this.database.addService(new Service(inputService.getName(), inputService.getUrlString()));

    }

    public void onAddServiceClick(ActionEvent actionEvent) {
        String name = serviceName.getText();
        String url = serviceUrl.getText();

        if (name.isEmpty() || url.isEmpty()) {
            alertContainer.getChildren().add(createAlert("Name or url should not be empty"));
        } else {

            Service service = new Service(name, url);
            addServiceToObservable(service);


            serviceName.setText("");
            serviceUrl.setText("");
            removeAlert();
        }
    }


    private void removeAlert() {
        alertContainer.getChildren().clear();

    }

    private HBox createAlert(String content) {
        HBox alert = new HBox();
        alert.getStyleClass().addAll("alert", "alert-danger");

        Text text = new Text();
        text.setText(content);

        alert.getChildren().add(text);

        return alert;
    }

    public void onCheckServiceClick(ActionEvent actionEvent) {
        this.manager.checkServicesStatus();

        ArrayList<Service> managerServices = this.manager.getServices();

        for (int i = 0; i < managerServices.size(); i++) {
            this.activeServices.set(i, managerServices.get(i));
        }

    }
}

